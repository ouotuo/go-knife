package bean

import(
	"reflect"
	"strconv"
	"fmt"
	"strings"
	"regexp"
)


const(
	TAG_BEAN="json"
	TAG_REQUIRE="require"
	TAG_REGEXP="regexp"
	TAG_DEFAULT="default"
)


//是否基本类型
func IsFieldPrimitive(field reflect.Value)(bool){
	switch field.Kind(){
	case reflect.Bool:
		return true
	case reflect.String:
		return true
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		return true
	case reflect.Float32,reflect.Float64:
		return true
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		return true

	default :
		return false
	}
}

//设置基本类型数值
func SetPrimitiveFieldStringValue(field reflect.Value,val string)(err error){
	switch field.Kind(){
	case reflect.Bool:
		var b bool
		b,err=strconv.ParseBool(val);if err==nil{
			field.SetBool(b)
		}
	case reflect.String:
		field.SetString(val)
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		var intVal int64
		intVal,err=strconv.ParseInt(val,10,0);if err==nil{
			field.SetInt(intVal)
		}
	case reflect.Float32,reflect.Float64:
		var floatVal float64
		floatVal,err=strconv.ParseFloat(val,0);if err==nil{
			field.SetFloat(floatVal)
		}
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		var uintVal uint64
		uintVal,err=strconv.ParseUint(val,10,0);if err==nil{
			field.SetUint(uintVal)
		}
	default :
		err=fmt.Errorf("not primitive type %v",field.Kind())
	}
	return
}

func SetBeanMap(pB interface{},vals map[string][]string)(err error){
	return setBeanMap(pB,"",vals)
}

//设置struct的属性
func setBeanMap(pB interface{},prefix string,vals map[string][]string)(err error){
	beanVal:=reflect.ValueOf(pB)

	if beanVal.Kind()!=reflect.Ptr {
		//必须是指针
		err=fmt.Errorf("pB argument %s should pointer",prefix)
		return
	}

	ele:=beanVal.Elem()

	if ele.Kind()!=reflect.Struct {
		return
	}

	//对所有field循环
	eleType:=ele.Type()
	var nfLen=ele.NumField()

	if prefix!=""{
		prefix=prefix+"."
	}

	for i:=0;i<nfLen;i++{
		sf:=eleType.Field(i)
		nf:=ele.Field(i)

		if nf.CanSet()==false{
			continue
		}

		var key string=sf.Tag.Get(TAG_BEAN)
		if key==""{
			key=fmt.Sprintf("%s%s",strings.ToLower(sf.Name[:1]),sf.Name[1:])
		}else if key=="-"{
			continue
		}

		key=prefix+key
		if IsFieldPrimitive(nf)==false{
			//不是基础类型
			if nf.Kind()==reflect.Ptr{
				if nf.IsNil(){
					//为nil，新建一个
					nf.Set(reflect.New(sf.Type.Elem()))
				}
				err=setBeanMap(nf.Interface(),key,vals)
			}
		}else{
			var strs=vals[key]
			//validate
			isRequired,_:=strconv.ParseBool(sf.Tag.Get(TAG_REQUIRE))

			var str=""
			if len(strs)>0{
				str=strs[0]
			}
			if str==""{
				//没有数据
				if isRequired==true{
					err=fmt.Errorf("%s require",key)
					return
				}
				var defaultVal=sf.Tag.Get(TAG_DEFAULT)
				if defaultVal!=""{
					str=defaultVal
				}
				if str==""{
					continue
				}
			}

			var checkRegexp string=sf.Tag.Get(TAG_REGEXP)
			if checkRegexp!=""{
				var r *regexp.Regexp
				r,err=regexp.Compile(checkRegexp)
				if err!=nil{
					err=fmt.Errorf("%s regexp wrong format",key)
					return
				}
				if r.MatchString(str)==false{
					err=fmt.Errorf("%s value not match regexp",prefix)
					return
				}
			}

			//convert to value
			err=SetPrimitiveFieldStringValue(nf,str)
		}

		//get error
		if err!=nil{
			return
		}
	}

	return
}

func ValidBean(pb interface{})(error){
	return validBean(pb,"")
}


//设置struct的属性
func validBean(pB interface{},prefix string)(err error){
	beanVal:=reflect.ValueOf(pB)

	if beanVal.Kind()!=reflect.Ptr {
		//必须是指针
		err=fmt.Errorf("pB argument should pointer")
		return
	}

	ele:=beanVal.Elem()

	if ele.Kind()!=reflect.Struct {
		return
	}

	//对所有field循环
	eleType:=ele.Type()
	var nfLen=ele.NumField()

	if prefix!=""{
		prefix=prefix+"."
	}

	for i:=0;i<nfLen;i++{
		sf:=eleType.Field(i)
		nf:=ele.Field(i)

		if nf.CanSet()==false{
			continue
		}

		var key string=sf.Tag.Get(TAG_BEAN)
		if key==""{
			key=fmt.Sprintf("%s%s",strings.ToLower(sf.Name[:1]),sf.Name[1:])
		}
		key=prefix+key
		if IsFieldPrimitive(nf)==false{
			//不是基础类型
			if nf.Kind()==reflect.Ptr{
				//是不是必须的

				if nf.IsNil(){
					//为nil，检查required
					isRequired,_:=strconv.ParseBool(sf.Tag.Get(TAG_REQUIRE))
					if isRequired==true{
						err=fmt.Errorf("%s is required",key)
						return
					}
				}else{
					err=validBean(nf.Interface(),key)
				}
			}
		}else{
			isRequired,_:=strconv.ParseBool(sf.Tag.Get(TAG_REQUIRE))
			var checkRegexp string=sf.Tag.Get(TAG_REGEXP)

			if isRequired==false && checkRegexp==""{
				continue
			}

			var str=fmt.Sprintf("%v",nf.Interface())
			//validate

			if str==""{
				//没有数据
				if isRequired==true{
					err=fmt.Errorf("%s require",key)
					return
				}
				continue
			}

			if checkRegexp!=""{
				var r *regexp.Regexp
				r,err=regexp.Compile(checkRegexp)
				if err!=nil{
					err=fmt.Errorf("%s regexp wrong format",key)
					return
				}
				if r.MatchString(str)==false{
					err=fmt.Errorf("%s value not match regexp",key)
					return
				}
			}
		}

		//get error
		if err!=nil{
			return
		}
	}

	return
}
