package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// SendPrinter send code to printer
func (env *Env) SendPrinter(printerName string) {
	log.Printf("Printer-%s start-up\n", printerName)
	errorHandler := func(id int64, err error) {
		log.Println(err)
		env.printQueue <- id
	}

	var id int64
	ok := true
	for ok {
		if id, ok = <-env.printQueue; ok {
			log.Println(id)
			p, err := env.db.GetPrintCode(id)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetPrintCode failed with %s", err))
				continue
			}
			u, err := env.db.GetUserAccount(p.Account)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetUserAccount failed with %s", err))
				break
			}
			log.Println(p)

			cmd := exec.Command("lp",
				"-d", printerName,
				"-t", u.DisplayName+"-"+u.SeatID,
				"-o", "prettyprint",
				"-o", "Page-left=36",
				"-o", "Page-right=36",
				"-o", "Page-top=36",
				"-o", "Page-bottom=36")
			cmd.Stdin = strings.NewReader(p.Code)
			err = cmd.Run()
			if err != nil {
				errorHandler(id, fmt.Errorf("Run cmd lp failed with %s", err))
				continue
			}
			p.IsDone = true
			err = env.db.UpdatePrintCode(*p)
			if err != nil {
				errorHandler(id, fmt.Errorf("Run cmd lp failed with %s", err))
			}
		}
	}
}

// PostPrinter append code to env.printQueue
func (env *Env) PostPrinter(c *gin.Context) {
	var requestJSON struct {
		PrintContent string `json:"PrintContent" binding:"required"`
	}
	err := c.BindJSON(&requestJSON)
	if err != nil {
		errMsg := fmt.Sprint("BindJSON failed with", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	raw, has := c.Get("user")
	if !has {
		c.AbortWithError(http.StatusInternalServerError, errors.New("No user in the gin.Context"))
		return
	}
	user := raw.(model.User)
	p, err := env.db.SavePrintCode(user.Account, requestJSON.PrintContent)
	if err != nil {
		errMsg := fmt.Sprint("SavePrintCode failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	log.Println("put", p)
	env.printQueue <- p.ID
	c.JSON(http.StatusOK, gin.H{"message": "OK", "queue_size": len(env.printQueue)})
}
