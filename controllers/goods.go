package controllers

import (
	"github.com/astaxie/beego/orm"
	"181105b/models"
	"math"
	"time"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type GoodsController struct {
	BaseController
}

func (this *GoodsController) ShowGoodsList()  {
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

	var goodsList []models.GoodsSKU
	s:=o.QueryTable("GoodsSKU")
	typename:=this.GetString("select")
	if typename=="" {
		s=s.RelatedSel("GoodsSPU").RelatedSel("GoodsSPU__GoodsType")
	} else {
		s=s.RelatedSel("GoodsSPU").RelatedSel("GoodsSPU__GoodsType").Filter("GoodsSPU__GoodsType__TypeName",typename)
	}

	totalRows,_:=s.Count()
	pageSize:=2
	pageCount:=int(math.Ceil(float64(totalRows)/float64(pageSize)))
	page,err:=this.GetInt("p")
	if err!=nil {
		page=1
	}
	if page>pageCount {
		page=pageCount
	}
	if page<1 {
		page=1
	}
	offset:=pageSize*(page-1)
	s.Limit(pageSize,offset).All(&goodsList)

	this.Data["page"]=page
	this.Data["totalRows"]=totalRows
	this.Data["pageCount"]=pageCount
	this.Data["goodsList"]=goodsList
	this.Data["typename"]=typename
	this.ShowLayout("首页")
	this.TplName="goods/index.html"
}

func (this *GoodsController) ShowAddGoodsSKU()  {
	goodsId,err:=this.GetInt("goods_id")
	if err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}
	var goods models.GoodsSPU
	goods.Id=goodsId
	o:=orm.NewOrm()
	if err=o.Read(&goods);err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}

	var goodsList []models.GoodsSKU
	s:=o.QueryTable("GoodsSKU").RelatedSel("GoodsSPU").Filter("GoodsSPU__Id",goodsId)
	totalRows,_:=s.Count()
	pageSize:=2
	pageCount:=int(math.Ceil(float64(totalRows)/float64(pageSize)))
	page,err:=this.GetInt("p")
	if err!=nil {
		page=1
	}
	if page>pageCount {
		page=pageCount
	}
	if page<1 {
		page=1
	}
	offset:=pageSize*(page-1)
	s.Limit(pageSize,offset).All(&goodsList)

	this.Data["goods"]=goods
	this.Data["page"]=page
	this.Data["totalRows"]=totalRows
	this.Data["pageCount"]=pageCount
	this.Data["goodsList"]=goodsList
	this.ShowLayout(goods.Name+" SKU")
	this.TplName="goods/add.html"
}

func (this *GoodsController) HandleAddGoodsSKU()  {
	goodsId,err:=this.GetInt("goods_id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goodsspu models.GoodsSPU
	goodsspu.Id=goodsId
	o:=orm.NewOrm()
	if err=o.Read(&goodsspu);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}

	name:=this.GetString("name")
	if name=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input name"}
		this.ServeJSON()
		return
	}

	stock,err:=this.GetInt("stock")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"input stock"}
		this.ServeJSON()
		return
	}
	if stock<0 {
		this.Data["json"]=&ResponseJSON{-1,nil,"invalid stock"}
		this.ServeJSON()
		return
	}

	price,err:=this.GetFloat("price")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"input price"}
		this.ServeJSON()
		return
	}
	if price<0 {
		this.Data["json"]=&ResponseJSON{-1,nil,"invalid price"}
		this.ServeJSON()
		return
	}

	unit:=this.GetString("unit")
	if unit=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input unit"}
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

	detail:=this.GetString("detail")
	var goods models.GoodsSKU
	goods.Name=name
	goods.Price=price
	goods.Unit=unit
	goods.Stock=stock
	goods.Image=image
	goods.Detail=detail
	loc,_:=time.LoadLocation("Asia/Saigon")
	goods.Addtime=time.Now().In(loc)
	goods.GoodsSPU=&goodsspu
	if _,err=o.Insert(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to add goodsSKU"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
	return
}

func (this *GoodsController) ShowGoodsSKUDetail()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/",302)
		return
	}
	var goods models.GoodsSKU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.QueryTable("GoodsSKU").RelatedSel("GoodsSPU").Filter("Id",id).One(&goods);err!=nil {
		this.Redirect("/",302)
		return
	}
	this.Data["goods"]=goods
	this.ShowLayout("商品SKU详情")
	this.TplName="goods/content.html"
}

func (this *GoodsController) DeleteGoodsSKU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goods models.GoodsSKU
	goods.Id=id
	o:=orm.NewOrm()
	if _,err=o.Delete(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to delete"}
		this.ServeJSON()
	} else {
		this.Data["json"]=&ResponseJSON{0,nil,"OK"}
		this.ServeJSON()
	}
}

func (this *GoodsController) ShowUpdateGoodsSKU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/",302)
		return
	}
	var goods models.GoodsSKU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.QueryTable("GoodsSKU").RelatedSel("GoodsSPU").Filter("Id",id).One(&goods);err!=nil {
		this.Redirect("/",302)
		return
	}
	this.Data["goods"]=goods
	this.ShowLayout("更新商品SKU")
	this.TplName="goods/update.html"
}

func (this *GoodsController) HandleUpdateGoodsSKU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goods models.GoodsSKU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.Read(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}

	name:=this.GetString("name")
	if name=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input name"}
		this.ServeJSON()
		return
	}

	stock,err:=this.GetInt("stock")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"input stock"}
		this.ServeJSON()
		return
	}
	if stock<0 {
		this.Data["json"]=&ResponseJSON{-1,nil,"invalid stock"}
		this.ServeJSON()
		return
	}

	price,err:=this.GetFloat("price")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"input price"}
		this.ServeJSON()
		return
	}
	if price<0 {
		this.Data["json"]=&ResponseJSON{-1,nil,"invalid price"}
		this.ServeJSON()
		return
	}

	unit:=this.GetString("unit")
	if unit=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"input unit"}
		this.ServeJSON()
		return
	}

	image,err:=this.UploadFile("image",[]string{".jpg",".png",".jpeg"},1024*1024*2)
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,err.Error()}
		this.ServeJSON()
		return
	}

	detail:=this.GetString("detail")

	goods.Name=name
	goods.Price=price
	goods.Unit=unit
	goods.Stock=stock
	if image!="" {
		goods.Image=image
	}
	goods.Detail=detail
	loc,_:=time.LoadLocation("Asia/Saigon")
	goods.Addtime=time.Now().In(loc)
	if _,err=o.Update(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to update goodsSKU"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
	return
}