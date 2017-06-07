package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os/exec"
	"strings"
)

// SendPrinter send code to printer
func (env *Env) SendPrinter(printerName string) {
	errorHandler := func(id int64, err error) {
		log.Println(err)
		env.printQueue <- id
	}

	var id int64
	ok := true
	for ok {
		if id, ok = <-env.printQueue; ok {
			p, err := env.db.GetPrintCode(id)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetPrintCode failed with %s", err))
				continue
			}
			u, err := env.db.GetUserAccount(p.Account)
			if err != nil {
				errorHandler(id, fmt.Errorf("GetUserAccount failed with %s", err))
				continue
			}

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
		c.JSON(400, gin.H{"message": errMsg})
		return
	}

	account, _ := c.Get("account")
	p, err := env.db.SavePrintCode(account.(string), requestJSON.PrintContent)
	if err != nil {
		errMsg := fmt.Sprint("SavePrintCode failed with", err)
		c.JSON(500, gin.H{"message": errMsg})
		return
	}
	env.printQueue <- p.ID
	c.JSON(200, gin.H{"message": "OK", "queue_size": len(env.printQueue)})
}
