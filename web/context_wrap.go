package web


const(
	HEADER_UA="User-Agent"
	HEADER_CONTEXT_TYPE="Context-Type"
)

type ReqContext interface {
	GetIp()string
	GetHeader(name string)string
	GetContentType()string
	Bind(form interface{})error
}


type ContextWrap struct{
	c ReqContext

	userId string
	userName string
}
func NewContextWrap(c ReqContext)(cw *ContextWrap){
	cw=&ContextWrap{
		c:c,
	}
	return
}
func(cw *ContextWrap)GetIp()string{
	return cw.c.GetIp()
}
func(cw *ContextWrap)GetHeader(name string)string{
	return cw.c.GetHeader(name)
}
func(cw *ContextWrap)GetUa()string{
	return cw.c.GetHeader(HEADER_UA)
}

func(cw *ContextWrap)SetUserId(userId string){
	cw.userId=userId
}
func(cw *ContextWrap)GetUserId()string{
	return cw.userId
}

func(cw *ContextWrap)SetUserName(userName string){
	cw.userName=userName
}
func(cw *ContextWrap)GetUserName()string{
	return cw.userName
}
func(cw *ContextWrap)GetContextType()string{
	return cw.c.GetContentType()
}

func(cw *ContextWrap)Bind(form interface{})error{
	return cw.c.Bind(form)
}


