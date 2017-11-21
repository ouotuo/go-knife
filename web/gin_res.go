package web

import (
	"reflect"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var _P_ERROR_TYPE=reflect.TypeOf((*error)(nil))
var _ERROR_TYPE=_P_ERROR_TYPE.Elem()
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
			log.Fatalf("%s,arguments[0] should ContextWrap",name)
		}
		needContextWrap=true
		var arg2=funType.In(1)
		if arg2.Kind()!=reflect.Ptr{
			log.Fatalf("%s,arguments[1] should ptr to struct ",name)
		}
		needBindForm=true
		formType=arg2
	}else{
		log.Fatalf("%s,arguments.len=%d > 2",name,numIn)
	}

	//check out
	if numOut==1{
		//(error)
		var out1=funType.Out(0)
		if out1!=_ERROR_TYPE{
			log.Fatalf("%s,out[0] should error",name)
		}
	}else if numOut==2{
		//(error,data)
		var out1=funType.Out(0)
		if out1!=_ERROR_TYPE{
			log.Fatalf("%s,out[0] should error",name)
		}
	}else{
		log.Fatalf("%s,out.len=%d,>2",name,numOut)
	}

	var funVal =reflect.ValueOf(fun)

	return func(c *gin.Context) {
		var result = &Result{}

		var arguments = make([]reflect.Value,numIn,numIn)
		var err error

		defer func() {
			if r := recover(); r != nil {
				log.Printf("execute exception,recover,%v", r)
				err,ok:=r.(error)
				if ok{
					result.SetErrSystemFail(err.Error())
				}else{
					result.SetErrSystemFail("error unknown type")
				}
			}
			c.JSON(200,result)
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
			//error
			if outs[0].IsNil()==false{
				result.SetErrParam(fmt.Sprintf("%v",outs[1].Interface()))
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
	if t.Kind()==reflect.Map{
		v=reflect.MakeMap(t)
	}else{
		v=reflect.New(t.Elem())
	}
	return
}