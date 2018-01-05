package web

import (
	"reflect"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/go-errors/errors"
)

var _P_ERROR_TYPE=reflect.TypeOf((*error)(nil))
var _ERROR_TYPE=_P_ERROR_TYPE.Elem()

var _P_IRESULT_INTERFACE=reflect.TypeOf((*IResult)(nil))
var _IRESULT_INTERFACE=_P_IRESULT_INTERFACE.Elem()

func CreateGinFunc(fun interface{})func(*gin.Context){
	var funType=reflect.TypeOf(fun)

	if funType.Kind()!=reflect.Func{
		log.Fatal("CreateGinFunc func argument is not function")
	}

	var needContextWrap,needBindForm bool
	var numIn=funType.NumIn()
	var numOut=funType.NumOut()

	var formType reflect.Type

	var name=funType.String()

	//check in
	if numIn==1{
		//(contextWrap) || (form)
		var arg1=funType.In(0)
		if arg1==reflect.TypeOf(&ContextWrap{}){
			needContextWrap=true
		}else if arg1.Kind()==reflect.Ptr{
			formType=arg1
			needBindForm=true
		}else{
			log.Fatalf("func[%s](ContextWrap|form), form argument should ptr to struct")
		}
	}else if numIn==2{
		//(contextWrap,form)
		var arg1=funType.In(0)
		if arg1!=reflect.TypeOf(&ContextWrap{}){
			log.Fatalf("%s,arguments[0] should *ContextWrap",name)
		}
		needContextWrap=true
		var arg2=funType.In(1)
		if arg2.Kind()==reflect.Ptr{
			needBindForm=true
			formType=arg2
		}else{
			log.Fatalf("%s,arguments[1] should ptr to struct ",name)
		}
	}else if numIn>2{
		log.Fatalf("%s,arguments.len=%d > 2",name,numIn)
	}

	//check out
	var isOutIResult bool
	var out1=funType.Out(0)
	if numOut==1{
		//(error)
		if out1==_ERROR_TYPE {
			//error类型
		}else if(out1.Implements(_IRESULT_INTERFACE)){
			isOutIResult=true
		}else{
			log.Fatalf("%s,out[0] should error",name)
		}
	}else if numOut==2{
		//(error,data)
		if out1!=_ERROR_TYPE{
			log.Fatalf("%s,out[0] should error",name)
		}
	}else{
		log.Fatalf("%s,out.len=%d,>2",name,numOut)
	}

	var funVal =reflect.ValueOf(fun)

	return func(c *gin.Context) {
		var result IResult
		if isOutIResult==false{
			result = &Result{}
		}else{
			result=reflect.New(out1)
		}

		var arguments = make([]reflect.Value,numIn,numIn)
		var err error

		defer func() {
			if r := recover(); r != nil {
				log.Printf("execute exception,recover,%v", r)
				e:=errors.New(r)
				log.Println(e.ErrorStack())
				err,ok:=r.(error)
				if ok{
					result.SetErrSystemFail(err.Error())
				}else{
					result.SetErrSystemFail("execute error,unknown type")
				}
				c.JSON(500,result)
			}else if result.IsOk(){
				c.JSON(200,result)
			}else if result.IsErrParam(){
				c.JSON(400,result)
			}else{
				c.JSON(205,result)
			}
		}()

		var index=0
		contextWrap:=NewGinContextWrap(c)
		if needContextWrap{
			arguments[index]=reflect.ValueOf(contextWrap)
			index++
		}
		if needBindForm{
			formValue:=createFormValue(formType)
			arguments[index]=formValue
			err=contextWrap.Bind(formValue.Interface())
			if err!=nil{
				result.SetErrParam(err.Error())
				return
			}
		}

		outs:= funVal.Call(arguments)

		result.SetOk(nil)

		if numOut>0{
			if isOutIResult{
				result=outs[0]
				return
			}

			//error
			if outs[0].IsNil()==false{
				result.SetErrExecute(fmt.Sprintf("%v",outs[0].Interface()))
				return
			}
		}
		if numOut>1{
			//data
			if outs[1].IsNil()==false{
				result.SetOk(outs[1].Interface())
			}
		}
	}
}

func createFormValue(t reflect.Type)(v reflect.Value){
	if t.Kind()==reflect.Map {
		v = reflect.MakeMap(t)
	}else if t.Kind()==reflect.Slice{
		v=reflect.MakeSlice(t,0,0)
	}else{
		v=reflect.New(t.Elem())
	}
	return
}