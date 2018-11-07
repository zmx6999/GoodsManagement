package routers

import (
	"181105b/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/",beego.BeforeExec,Filter)
	beego.InsertFilter("/goods/*",beego.BeforeExec,Filter)
	beego.InsertFilter("/type/*",beego.BeforeExec,Filter)
	beego.InsertFilter("/goodsspu/*",beego.BeforeExec,Filter)

	beego.Router("/",&controllers.GoodsController{},"get:ShowGoodsList")

	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")

	beego.Router("/type/add",&controllers.GoodsTypeController{},"get:ShowType;post:HandleAddType")
	beego.Router("/type/delete",&controllers.GoodsTypeController{},"post:DeleteType")
	beego.Router("/type/update",&controllers.GoodsTypeController{},"get:ShowUpdateType;post:HandleUpdateType")
	beego.Router("/type/detail",&controllers.GoodsTypeController{},"get:ShowTypeDetail")

	beego.Router("/goodsspu/add",&controllers.GoodsSPUController{},"get:ShowGoodsSPU;post:HandleAddGoodsSPU")
	beego.Router("/goodsspu/delete",&controllers.GoodsSPUController{},"post:DeleteGoodsSPU")
	beego.Router("/goodsspu/detail",&controllers.GoodsSPUController{},"get:ShowGoodsSPUDetail")
	beego.Router("/goodsspu/update",&controllers.GoodsSPUController{},"get:ShowUpdateGoodsSPU;post:HandleUpdateGoodsSPU")

	beego.Router("/goods/add",&controllers.GoodsController{},"get:ShowAddGoodsSKU;post:HandleAddGoodsSKU")
	beego.Router("/goods/detail",&controllers.GoodsController{},"get:ShowGoodsSKUDetail")
	beego.Router("/goods/delete",&controllers.GoodsController{},"post:DeleteGoodsSKU")
	beego.Router("/goods/update",&controllers.GoodsController{},"get:ShowUpdateGoodsSKU;post:HandleUpdateGoodsSKU")
}

var Filter = func(ctx *context.Context) {
	username:=ctx.Input.Session("username")
	if username==nil {
		ctx.Redirect(302,"/login")
	}
}