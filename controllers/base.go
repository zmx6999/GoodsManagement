package controllers

import (
	"github.com/astaxie/beego"
	"crypto/sha256"
	"encoding/hex"
	"path"
	"errors"
	"strings"
	"strconv"
	"github.com/weilaihui/fdfs_client"
	"net/http"
)

type BaseController struct {
	beego.Controller
}

type ResponseJSON struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

func (this *BaseController) Sha256Str(x string) string {
	y:=sha256.Sum256([]byte(x))
	return hex.EncodeToString(y[:])
}

func (this *BaseController) Find(x string,a []string) int {
	for k,v:=range a{
		if x==v {
			return k
		}
	}
	return -1
}

func (this *BaseController) UploadFile(from string,tps []string,maxsize int64) (string,error) {
	file,head,err:=this.GetFile(from)
	if err!=nil {
		if err==http.ErrMissingFile {
			return "",nil
		} else {
			return "",err
		}
	}
	defer file.Close()
	if head.Filename=="" {
		return "",nil
	}
	ext:=strings.ToLower(path.Ext(head.Filename))
	if this.Find(ext,tps)<0 {
		return "",errors.New(from+" type can only be "+strings.Join(tps,","))
	}
	if head.Size>maxsize {
		return "",errors.New(from+"<"+strconv.Itoa(int(maxsize))+"B")
	}
	m:=make([]byte,head.Size)
	file.Read(m)
	client,err:=fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err!=nil {
		return "",err
	}
	rsp,err:=client.UploadByBuffer(m,"")
	if err!=nil {
		return "",err
	}
	return rsp.RemoteFileId,nil
}

func (this *BaseController) ShowLayout(title string)  {
	username:=this.GetSession("username")
	this.Data["username"]=username.(string)
	this.Data["title"]=title
	this.Layout="layout.html"
}