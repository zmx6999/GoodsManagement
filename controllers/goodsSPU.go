package controllers

import (
	"181105b/models"
	"github.com/astaxie/beego/orm"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type GoodsSPUController struct {
	BaseController
}

func (this *GoodsSPUController) ShowGoodsSPU()  {
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

	var goodsList []models.GoodsSPU
	s:=o.QueryTable("GoodsSPU")
	typename:=this.GetString("select")
	if typename=="" {
		s=s.RelatedSel("GoodsType")
	} else {
		s=s.RelatedSel("GoodsType").Filter("GoodsType__TypeName",typename)
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
	this.ShowLayout("商品")
	this.TplName="goodsSPU/add.html"

}

func (this *GoodsSPUController) HandleAddGoodsSPU()  {
	name:=this.GetString("name")
	typename:=this.GetString("typename")
	if name=="" || typename=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"imcomplete"}
		this.ServeJSON()
		return
	}
	desc:=this.GetString("desc")
	var goodsType models.GoodsType
	goodsType.TypeName=typename
	o:=orm.NewOrm()
	o.Read(&goodsType,"TypeName")

	var goods models.GoodsSPU
	goods.Name=name
	goods.Desc=desc
	goods.GoodsType=&goodsType
	if _,err:=o.Insert(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to add goods"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
	return
}

func (this *GoodsSPUController) DeleteGoodsSPU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goods models.GoodsSPU
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

func (this *GoodsSPUController) ShowGoodsSPUDetail()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}
	var goods models.GoodsSPU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.QueryTable("GoodsSPU").RelatedSel("GoodsType").Filter("Id",id).One(&goods);err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}
	this.Data["goods"]=goods
	this.ShowLayout("商品SPU详情")
	this.TplName="goodsSPU/content.html"
}

func (this *GoodsSPUController) ShowUpdateGoodsSPU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}
	var goods models.GoodsSPU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.QueryTable("GoodsSPU").RelatedSel("GoodsType").Filter("Id",id).One(&goods);err!=nil {
		this.Redirect("/goodsspu/add",302)
		return
	}
	this.Data["goods"]=goods

	var types []models.GoodsType
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

	this.ShowLayout("更新商品SPU")
	this.TplName="goodsSPU/update.html"
}

func (this *GoodsSPUController) HandleUpdateGoodsSPU()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}
	var goods models.GoodsSPU
	goods.Id=id
	o:=orm.NewOrm()
	if err=o.Read(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"not found"}
		this.ServeJSON()
		return
	}

	name:=this.GetString("name")
	typename:=this.GetString("typename")
	if name=="" || typename=="" {
		this.Data["json"]=&ResponseJSON{-1,nil,"imcomplete"}
		this.ServeJSON()
		return
	}
	desc:=this.GetString("desc")

	var goodsType models.GoodsType
	goodsType.TypeName=typename
	o.Read(&goodsType,"TypeName")

	goods.Name=name
	goods.Desc=desc
	goods.GoodsType=&goodsType
	if _,err:=o.Update(&goods);err!=nil {
		this.Data["json"]=&ResponseJSON{-1,nil,"failed to update goods"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=&ResponseJSON{0,nil,"OK"}
	this.ServeJSON()
	return
}