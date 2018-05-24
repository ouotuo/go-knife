package cfg

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/ouotuo/go-knife/bean"
	"github.com/sakeven/go-env"
	log "github.com/sirupsen/logrus"
)

const DEFAULT_FILE  = "cfg.json"

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

//file or env
func Load(file string,ptr interface{})(err error){
	if _,err:=os.Stat(file);err!=nil{
		//load from env
		log.Infof("load from env,%s",os.Environ())
		err=env.Decode(ptr)
	}else{
		log.Infof("load from file,%s",file)
		err=LoadFromJsonFile(file,ptr)
	}
	return
}