package web

type Result struct{
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
	Data interface{} `json:"data"`
}

type IResult interface {
	IsOk()bool
	SetOk(data interface{})
	IsErrParam()bool
	SetErrParam(msg string)
	IsErrExecute()bool
	SetErrExecute(msg string)
	IsErrSystemFail()bool
	SetErrSystemFail(msg string)
}

const(
	ERRCODE_OK=0
	ERRCODE_PARAM=400
	ERRCODE_EXECUTE=205
	ERRCODE_SYSTEM_FAIL=-1
)


func(r *Result)SetResult(errcode int,errmsg string)(*Result){
	r.Errcode=errcode
	r.Errmsg=errmsg
	return r
}
func(r *Result)SetErrParam(msg string){
	r.SetResult(ERRCODE_PARAM,msg)
}

func(r *Result)SetErrExecute(msg string){
	r.SetResult(ERRCODE_EXECUTE,msg)
}

func(r *Result)SetErrSystemFail(msg string){
	r.SetResult(ERRCODE_SYSTEM_FAIL,msg)
}

func(r *Result)SetOk(data interface{}){
	r.SetResult(ERRCODE_OK,"ok")
	r.Data=data
}

func(r *Result)IsOk()bool{
	return r.Errcode==ERRCODE_OK
}

func(r *Result)IsErrParam()bool{
	return r.Errcode==ERRCODE_PARAM
}

func(r *Result)IsErrExecute()bool{
	return r.Errcode==ERRCODE_EXECUTE
}
func(r *Result)IsErrSystemFail()bool{
	return r.Errcode==ERRCODE_SYSTEM_FAIL
}

