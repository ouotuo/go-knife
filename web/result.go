package web

type Result struct{
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
	Data interface{} `json:"data"`
}

const(
	ERRCODE_OK=0
	ERRCODE_PARAM=400
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

func(r *Result)SetErrSystemFail(msg string){
	r.SetResult(ERRCODE_SYSTEM_FAIL,msg)
}

func(r *Result)SetOk(data interface{}){
	r.SetResult(ERRCODE_OK,"ok")
	r.Data=data
}

