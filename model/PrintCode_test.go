package model

import (
	"log"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Test_PrintCode(t *testing.T) {
	db, err := OpenDBTest()
	if err != nil {
		t.Error(err)
	}
	p, err := db.SavePrintCode("admin", `
	#include <iostream>

	using namespace std;

	int main() {
		std::cout<<"Hello World!"<<std::endl;
		return 0;
	}
	`)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Duration(1) * time.Second)
	log.Println(p)
	p, err = db.SavePrintCode("ZhuangZhou", `
	print "鲲之大"
	`)
	if err != nil {
		t.Error(err)
	}
	log.Println(p)

	p.Code = `print "一锅炖不下"`
	err = db.UpdatePrintCode(*p)
	if err != nil {
		t.Error(err)
	}
	p, err = db.GetPrintCode(p.ID)
	if err != nil {
		t.Error(err)
	}
	log.Println(p)
}
