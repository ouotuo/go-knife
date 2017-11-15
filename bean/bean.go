package bean

import(
	"reflect"
	"strconv"
)


func SetFieldStringValue(field reflect.Value,val string)(err error){
	switch(field.Kind()){
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
		err=fmt.Errorf("unknown type %v",field.Kind())
	}
	return
}
