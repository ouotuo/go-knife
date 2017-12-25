package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

const(
	FIELD_NAME_ID="Id"
)

type AddFun func(session *mgo.Session,bean interface{})(bson.ObjectId,error)

func GetAddFun(table string)AddFun{
	return func(session *mgo.Session,bean interface{})(retId bson.ObjectId,err error){
		c := session.DB("").C(table)

		//add objectId
		elem:=reflect.ValueOf(bean).Elem()
		field:=elem.FieldByName(FIELD_NAME_ID);if field.IsValid(){
			bsonId,ok:=field.Interface().(bson.ObjectId);if ok==true{
				if bsonId.Valid()==false{
					retId=bson.NewObjectId()
					field.Set(reflect.ValueOf(retId))
				}else{
					retId=bsonId
				}
			}
		}

		err= c.Insert(bean)
		return
	}
}

type BatchAddFun func(session *mgo.Session,beans []interface{})(error)

func GetBatchAddFun(table string)BatchAddFun{
	return func(session *mgo.Session,beans []interface{})(err error){
		c := session.DB("").C(table)

		return c.Insert(beans...)
	}
}


type UpdateFun func(session *mgo.Session,id bson.ObjectId,update bson.M)(error)
func GetUpdateFun(table string)UpdateFun{
	return func(session *mgo.Session,id bson.ObjectId,update bson.M)(error){
		c := session.DB("").C(table)

		return c.UpdateId(id,update)
	}
}

type DelFun func(session *mgo.Session,id bson.ObjectId)(error)
func GetDelFun(table string)DelFun{
	return func(session *mgo.Session,id bson.ObjectId)(error){
		c := session.DB("").C(table)

		return c.RemoveId(id)
	}
}

//返回指针
type NewInstanceFun func()(interface{})

type LoadFun func(session *mgo.Session,id bson.ObjectId)(interface{},error)
func GetLoadFun(table string,newFun NewInstanceFun)LoadFun{
	return func(session *mgo.Session,id bson.ObjectId)(interface{},error){
		c := session.DB("").C(table)

		var val interface{}=newFun()

		err:=c.FindId(id).One(val)
		if err==nil{
			return val,nil
		}else if err==mgo.ErrNotFound{
			return nil,nil
		}else{
			return nil,err
		}
	}
}


type NewInstanceSliceFun func()(interface{})
type ListFun func(session *mgo.Session,query bson.M,skip int,limit int,sort ...string)(interface{},error)
func GetListFun(table string,newFun NewInstanceSliceFun)ListFun{
	return func(session *mgo.Session,query bson.M,skip int,limit int,sort ...string)(interface{},error){
		c := session.DB("").C(table)

		var val interface{}=newFun()

		q:=c.Find(query)
		if skip>0{
			q=q.Skip(skip)
		}
		if limit>0{
			q=q.Limit(limit)
		}
		if len(sort)>0{
			q=q.Sort(sort...)
		}
		err:=q.All(val)
		if err==nil{
			return val,nil
		}else if err==mgo.ErrNotFound{
			return nil,nil
		}else{
			return nil,err
		}
	}
}

type PageFun func(session *mgo.Session,direction string,curId string,query bson.M,pageSize int)(*PageBean,error)
func GetPageFun(table string,newFun NewInstanceSliceFun)PageFun {
	return func(session *mgo.Session, direction string, curId string, query bson.M, pageSize int) (*PageBean, error) {
		c := session.DB("").C(table)
		var val interface{}=newFun()

		if pageSize<=0{
			pageSize=10
		}

		for k,v:=range query{
			if v==nil{
				delete(query,k)
				continue
			}
			str,ok:=v.(string);if ok && str==""{
				delete(query,k)
				continue
			}
		}
		return RevPage(c,query,direction,curId,true,pageSize,val)
	}
}

func GetRevPageFun(table string,newFun NewInstanceSliceFun)PageFun {
	return func(session *mgo.Session, direction string, curId string, query bson.M, pageSize int) (*PageBean, error) {
		c := session.DB("").C(table)
		var val interface{}=newFun()

		if pageSize<=0{
			pageSize=10
		}

		for k,v:=range query{
			if v==nil{
				delete(query,k)
				continue
			}
			str,ok:=v.(string);if ok && str==""{
				delete(query,k)
				continue
			}
		}
		return Page(c,query,direction,curId,true,pageSize,val)
	}
}

//count
type CountFun func(session *mgo.Session,query bson.M)(int,error)
func GetCountFun(table string)CountFun{
	return func(session *mgo.Session,query bson.M)(int,error){
		c := session.DB("").C(table)

		return c.Find(query).Count()
	}
}

