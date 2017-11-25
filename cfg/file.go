package cfg

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/ouotuo/go-knife/bean"
)

//从文件加载
func LoadFromJsonFile(file string,ptr interface{})(err error){
	f,err:=os.Open(file)
	if err!=nil{
		return
	}
	defer f.Close()

	bs,err:=ioutil.ReadAll(f)
	if err!=nil{
		return
	}

	err=json.Unmarshal(bs,ptr)

	return
}


func LoadFromJsonFileAndValid(file string,ptr interface{})(err error){
	err=LoadFromJsonFile(file,ptr)
	if err==nil{
		err=bean.ValidBean(ptr)
	}
	return
}

func ToJson(ptr interface{})string{
	bs,_:=json.Marshal(ptr)
	return string(bs)
}
