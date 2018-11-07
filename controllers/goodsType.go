package controllers

import (
	"181105b/models"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type GoodsTypeController struct {
	BaseController
}

func (this *GoodsTypeController) ShowType()  {
	var types []models.GoodsType
	o:=orm.NewOrm()
	conn,_:=redis.Dial("tcp","127.0.0.1:6379")
	tb,_:=redis.Bytes(conn.Do("get","types"))
	decoder:=gob.NewDecoder(bytes.NewReader(tb))
	decoder.Decode(&types)
	if len(types)==0 {
		o.QueryTable("GoodsType").All(&types)
		var buffer bytes.Buffer
		encoder:=gob.NewEncoder(&buffer)
		encoder.Encode(&types)
		conn.Do("set","types",buffer.Bytes())
	}
	this.Data["types"]=types
	this.ShowLayout("分类")
	this.TplName="goodsType/add.html"
}

func (this *GoodsTypeController) HandleAddType()  {
	typename:=this.GetString("typename")
	if typename=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input typename"}
		this.ServeJSON()
		return
	}
	logo,err:=this.UploadFile("logo",[]string{".jpg",".png",".jpeg"},1024*1024*2)
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,err.Error()}
		this.ServeJSON()
		return
	}
	if logo=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"upload logo"}
		this.ServeJSON()
		return
	}
	image,err:=this.UploadFile("image",[]string{".jpg",".png",".jpeg"},1024*1024*2)
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,err.Error()}
		this.ServeJSON()
		return
	}
	if image=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"upload image"}
		this.ServeJSON()
		return
	}
	var goodsType models.GoodsType
	goodsType.TypeName=typename
	goodsType.Logo=logo
	goodsType.Image=image
	o:=orm.NewOrm()
	if _,err=o.Insert(&goodsType);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to add type"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
	return
}

func (this *GoodsTypeController) DeleteType()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goodsType models.GoodsType
	goodsType.Id=id
	o:=orm.NewOrm()
	if _,err=o.Delete(&goodsType);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to delete"}
		this.ServeJSON()
	} else {
		this.Data["json"]=&ResponseJSON{0,nil,"OK"}
		this.ServeJSON()
	}
}

func (this *GoodsTypeController) ShowTypeDetail()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/type/add",302)
		return
	}
	var goodsType models.GoodsType
	goodsType.Id=id
	o:=orm.NewOrm()
	if err=o.Read(&goodsType);err!=nil {
		this.Redirect("/type/add",302)
		return
	}
	this.Data["goodsType"]=goodsType
	this.ShowLayout("分类详情")
	this.TplName="goodsType/content.html"
}

func (this *BaseController) ShowUpdateType()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/type/add",302)
		return
	}
	var goodsType models.GoodsType
	goodsType.Id=id
	o:=orm.NewOrm()
	if err=o.Read(&goodsType);err!=nil {
		this.Redirect("/type/add",302)
		return
	}
	this.Data["goodsType"]=goodsType
	this.ShowLayout("更新分类")
	this.TplName="goodsType/update.html"
}

func (this *BaseController) HandleUpdateType()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}

	var goodsType models.GoodsType
	goodsType.Id=id
	o:=orm.NewOrm()
	if err=o.Read(&goodsType);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}

	typename:=this.GetString("typename")
	if typename=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input typename"}
		this.ServeJSON()
		return
	}
	logo,err:=this.UploadFile("logo",[]string{".jpg",".png",".jpeg"},1024*1024*2)
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,err.Error()}
		this.ServeJSON()
		return
	}
	image,err:=this.UploadFile("image",[]string{".jpg",".png",".jpeg"},1024*1024*2)
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,err.Error()}
		this.ServeJSON()
		return
	}

	goodsType.TypeName=typename
	if logo!="" {
		goodsType.Logo=logo
	}
	if image!="" {
		goodsType.Image=image
	}
	if _,err:=o.Update(&goodsType);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to update type"}
		this.ServeJSON()
	}

	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
}