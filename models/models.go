package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	)

type User struct {
	Id int
	Username string `orm:"size(20);unique"`
	Password string `orm:"size(64)"`
	Email string `orm:"size(50)"`
	Active bool `orm:"default(false)"`
	UserType int `orm:"default(0)"`

	Addresses []*Address `orm:"reverse(many)"`
	Orders []*Order `orm:"reverse(many)"`
}

type Address struct {
	Id int
	Receiver string `orm:"size(20)"`
	Addr string `orm:"size(100)"`
	Phone string `orm:"size(20)"`
	IsDefault bool `orm:"default(false)"`

	User *User `orm:"rel(fk)"`

	Orders []*Order `orm:"reverse(many)"`
}

type GoodsSPU struct {
	Id int
	Name string `orm:"size(30)"`
	Desc string `orm:"size(100)"`

	GoodsType *GoodsType `orm:"rel(fk)"`

	GoodsSKUs []*GoodsSKU `orm:"reverse(many)"`
}

type GoodsSKU struct {
	Id int
	Name string `orm:"size(30)"`
	Price float64
	Unit string `orm:"size(20)"`
	Image string `orm:"size(100)"`
	Stock int `orm:"default(1)"`
	Sold int `orm:"default(0)"`
	Detail string `orm:"size(600)"`
	Addtime time.Time
	Status int `orm:"default(1)"`

	GoodsSPU *GoodsSPU `orm:"rel(fk)"`

	GoodsImages []*GoodsImage `orm:"reverse(many)"`
	GoodsIndexBanners []*GoodsIndexBanner `orm:"reverse(many)"`
	GoodsIndexTypes []*GoodsIndexType `orm:"reverse(many)"`
	OrderGoodss []*OrderGoods `orm:"reverse(many)"`
}

type GoodsType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Image string `orm:"size(100)"`
	Logo string `orm:"size(100)"`

	GoodsSPUs []*GoodsSPU `orm:"reverse(many)"`
	GoodsIndexTypes []*GoodsIndexType `orm:"reverse(many)"`
}

type GoodsImage struct {
	Id int
	Image string `orm:"size(100)"`

	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
}

type GoodsIndexBanner struct {
	Id int
	Image string `orm:"size(100)"`
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
	Index int `orm:"default(0)"`
}

type GoodsIndexActicity struct {
	Id int
	Name string `orm:"size(30)"`
	Image string `orm:"size(100)"`
	Url string `orm:"size(100)"`
	Index int `orm:"default(0)"`
}

type GoodsIndexType struct {
	Id int
	ShowType int `orm:"default(0)"`
	Index int `orm:"default(0)"`

	GoodsType *GoodsType `orm:"rel(fk)"`
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
}

type Order struct {
	Id int
	OrderId string `orm:"size(20)"`
	GoodsFee float64 `orm:"default(0)"`
	TransitFee float64 `orm:"default(0)"`
	PayMethod string `orm:"size(20)"`
	PayStatus int `orm:"default(0)"`
	PayId string `orm:"size(20)"`
	AddTime time.Time

	User *User `orm:"rel(fk)"`
	Address *Address `orm:"rel(fk)"`

	OrderGoodss []*OrderGoods `orm:"reverse(many)"`
}

type OrderGoods struct {
	Id int
	Price float64
	Amount int
	Comment string `orm:"size(200)"`

	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
	Order *Order `orm:"rel(fk)"`
}

func init()  {
	orm.RegisterDataBase("default","mysql","root:123456@tcp(localhost:3306)/db181105?charset=utf8")
	orm.RegisterModel(new(User),new(Address),new(GoodsSPU),new(GoodsSKU),new(GoodsType),new(GoodsImage),new(GoodsIndexBanner),new(GoodsIndexActicity),new(GoodsIndexType),new(Order),new(OrderGoods))
	orm.RunSyncdb("default",false,true)
}