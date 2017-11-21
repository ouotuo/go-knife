package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/ouotuo/go-knife/bean"
)

/***** context ****/
type ginReqContext  struct{
	c *gin.Context
}

func NewGinContextWrap(c *gin.Context)(cw *ContextWrap){
	cw=&ContextWrap{
		c:&ginReqContext{c:c},
	}
	return
}

func (gc *ginReqContext)GetIp()string{
	return gc.c.ClientIP()
}
func(gc *ginReqContext)GetContentType()string{
	return gc.c.ContentType()
}
func( gc *ginReqContext)GetHeader(name string)string{
	return gc.c.GetHeader(name)
}
func( gc *ginReqContext)Bind(form interface{})(err error){
	contentType:=gc.c.ContentType()

	switch contentType {
	case gin.MIMEJSON:
		err=ginBindJson(gc.c,form)
		if err==nil{
			err=bean.ValidBean(form)
		}
		return
	case gin.MIMEPOSTForm,gin.MIMEMultipartPOSTForm:
	default:
	}

	err =gc.c.Request.ParseForm()
	if err!=nil{
		return
	}

	var mapParams=gc.c.Request.Form
	if len(gc.c.Params)>0{
		if mapParams==nil{
			mapParams=make(map[string][]string)
		}
		for _,param:=range gc.c.Params{
			mapParams[param.Key]=[]string{param.Value}
		}
	}

	err=bean.SetBeanMap(form,mapParams)

	return
}

func ginBindJson(c *gin.Context,form interface{})(err error){
	body,err:=ioutil.ReadAll(c.Request.Body)

	if err!=nil{
		err=fmt.Errorf("read request body error,%s",err)
	}else{
		err=json.Unmarshal(body,form)
	}

	return
}

