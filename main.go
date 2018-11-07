package main

import (
	_ "181105b/routers"
	"github.com/astaxie/beego"
	"time"
)

func main() {
	beego.AddFuncMap("getconfig",GetConfig)
	beego.AddFuncMap("nextpage",ShowNextPage)
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("showtime",ShowTime)
	beego.Run()
}

func GetConfig(key string) string {
	return beego.AppConfig.String(key)
}

func ShowNextPage(x int,pageCount int) int {
	if x==pageCount {
		return pageCount
	} else {
		return x+1
	}
}

func ShowPrePage(x int) int {
	if x==1 {
		return 1
	} else {
		return x-1
	}
}

func ShowTime(t time.Time,timezone string) string {
	loc,_:=time.LoadLocation(timezone)
	return t.In(loc).Format("2006-01-02 15:04:05")
}