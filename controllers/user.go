package controllers

import (
	"github.com/astaxie/beego/orm"
	"181105b/models"
	"encoding/base64"
)

type UserController struct {
	BaseController
}

func (this *UserController) ShowRegister()  {
	this.TplName="user/register.html"
}

func (this *UserController) HandleRegister()  {
	username:=this.GetString("username")
	password:=this.GetString("password")
	if username=="" || password=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"imcomplete"}
		this.ServeJSON()
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Username=username
	user.Password=this.Sha256Str(password)
	if _,err:=o.Insert(&user);err!=nil{
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to register"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
}

func (this *UserController) ShowLogin()  {
	username:=this.Ctx.GetCookie("username")
	if username!="" {
		tmp,_:=base64.StdEncoding.DecodeString(username)
		this.Data["username"]=string(tmp)
		this.Data["remember"]="checked"
	}
	this.TplName="user/login.html"
}

func (this *UserController) HandleLogin()  {
	username:=this.GetString("username")
	password:=this.GetString("password")
	if username=="" || password=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"imcomplete"}
		this.ServeJSON()
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Username=username
	if err:=o.Read(&user,"Username");err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"error username"}
		this.ServeJSON()
		return
	}
	if user.Password!=this.Sha256Str(password) {
		this.Data["json"]=&ResponseJSON{-1,nil,"error password"}
		this.ServeJSON()
		return
	}
	remember:=this.GetString("remember")
	if remember!="" {
		this.Ctx.SetCookie("username",base64.StdEncoding.EncodeToString([]byte(username)),100)
	} else {
		this.Ctx.SetCookie("username","",-1)
	}
	this.SetSession("username",username)
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
}

func (this *UserController) Logout() {
	this.DelSession("username")
	this.Redirect("/login",302)
}