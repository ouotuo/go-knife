package web

import "github.com/gin-gonic/gin"

//requestWrap info


type ReqWrap interface{
	//info
	GetIp()string;
	GetHeader(name string)string;
	GetUa()string;

	//user
	GetUserId()int64;
	GetUserName()string;

}

const(
	HEADER_UA="User-Agent"
)

type GinReqWrap struct{
	c *gin.Context
	userId int64
	userName string
}

func NewGinReqWrap(c *gin.Context,userId int64,userName string)(grw *GinReqWrap){
	grw=&GinReqWrap{
		c:c,
		userId:userId,
		userName:userName,
	}
	return
}

func(grw *GinReqWrap)GetIp(){
	return grw.c.ClientIP()
}
func(grw *GinReqWrap)GetHeader(name string)string{
	return grw.c.GetHeader(name)
}
func(grw *GinReqWrap)GetUa()string{
	return grw.GetHeader(HEADER_UA)
}

func(grw *GinReqWrap)GetUserId()int64{
	return grw.userId
}
func(grw *GinReqWrap)GetUserName()string{
	return grw.userName
}
