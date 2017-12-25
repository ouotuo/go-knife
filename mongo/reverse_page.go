package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"math"
	"fmt"
)



func rev_beginPage(c *mgo.Collection,query map[string]interface{},pageSize int,result interface{})(*PageBean,error){
	pageBean:=&PageBean{Total:0,HasPrev:false,HasNext:false,PrevId:"",NextId:"",PageSize:pageSize,}

	limit:=pageSize+1

	sortField:="-_id"
	err:=c.Find(query).Sort(sortField).Limit(limit).All(result);if err!=nil {
		return nil,err
	}else{
		v:=reflect.ValueOf(result)
		ele:=v.Elem()
		if ele.Len()==limit{
			//说明还有下一页数据
			pageBean.HasNext=true
			lastEle:=ele.Index(pageSize-1).Elem()
			lastObjectId:=bson.ObjectId(lastEle.FieldByName("Id").String())
			pageBean.NextId=lastObjectId.Hex()
			ele.Set(ele.Slice(0,limit-1))
		}

		pageBean.List=result
		return pageBean,nil
	}
}



func rev_endPage(c *mgo.Collection,query map[string]interface{},pageSize int,result interface{})(*PageBean,error){
	pageBean:=&PageBean{Total:0,HasPrev:false,HasNext:false,PrevId:"",NextId:"",PageSize:pageSize,}

	limit:=pageSize+1

	sortField:="_id"

	err:=c.Find(query).Sort(sortField).Limit(limit).All(result);if err!=nil {
		return nil,err
	}else{
		v:=reflect.ValueOf(result)
		ele:=v.Elem()
		if ele.Len()==limit{
			//说明还有上一页数据
			pageBean.HasPrev=true
			lastEle:=ele.Index(limit-1).Elem()
			lastObjectId:=bson.ObjectId(lastEle.FieldByName("Id").String())
			pageBean.PrevId=lastObjectId.Hex()
			ele.Set(ele.Slice(0,limit-1))
		}
		pageBean.List=result
		reverseSlice(result)
		return pageBean,nil
	}
}

//前一个包括,后一个不包括
func rev_nextPage(c *mgo.Collection,query map[string]interface{},pageSize int,curId string,result interface{})(*PageBean,error){
	pageBean:=&PageBean{Total:0,HasPrev:true,HasNext:false,PrevId:curId,NextId:curId,PageSize:pageSize,}

	limit:=pageSize+1

	sortField:="-_id"

	if curId!="" && bson.IsObjectIdHex(curId)==false{
		return nil,fmt.Errorf("curId=%s is not hex objectId",curId)
	}

	curObjId:=bson.ObjectIdHex(curId)
	if query!=nil{
		query["_id"]=bson.M{"$lt":curObjId}
	}else{
		query=bson.M{"_id":bson.M{"$lt":curObjId}}
	}

	err:=c.Find(query).Sort(sortField).Limit(limit).All(result);if err!=nil {
		return nil,err
	}else{
		v:=reflect.ValueOf(result)
		ele:=v.Elem()
		if ele.Len()==limit{
			//说明还有下一页数据
			pageBean.HasNext=true
			lastEle:=ele.Index(pageSize-1).Elem()
			lastObjectId:=bson.ObjectId(lastEle.FieldByName("Id").String())
			pageBean.NextId=lastObjectId.Hex()
			ele.Set(ele.Slice(0,limit-1))
		}

		pageBean.List=result
		return pageBean,nil
	}
}

func rev_prevPage(c *mgo.Collection,query map[string]interface{},pageSize int,curId string,result interface{})(*PageBean,error){
	pageBean:=&PageBean{Total:0,HasPrev:false,HasNext:true,PrevId:"",NextId:curId,PageSize:pageSize,}

	limit:=pageSize+1

	sortField:="_id"

	curObjId:=bson.ObjectIdHex(curId)
	if query!=nil{
		query["_id"]=bson.M{"$gte":curObjId}
	}else{
		query=bson.M{"_id":bson.M{"$gte":curObjId}}
	}

	err:=c.Find(query).Sort(sortField).Limit(limit).All(result);if err!=nil {
		return nil,err
	}else{
		v:=reflect.ValueOf(result)
		ele:=v.Elem()
		if ele.Len()==limit{
			//说明还有上一页数据
			pageBean.HasPrev=true
			lastEle:=ele.Index(limit-1).Elem()
			lastObjectId:=bson.ObjectId(lastEle.FieldByName("Id").String())
			pageBean.PrevId=lastObjectId.Hex()
			ele.Set(ele.Slice(0,limit-1))
		}
		pageBean.List=result
		reverseSlice(result)
		return pageBean,nil
	}
}


func rev_queryPageBean(c *mgo.Collection,query map[string]interface{},direction string,pageSize int,curId string,result interface{})(*PageBean,error){
	isNext:=true
	if direction=="prev"{
		isNext=false
	}
	if isNext==true{
		if curId!=""{
			return rev_nextPage(c,query,pageSize,curId,result)
		}else{
			return rev_beginPage(c,query,pageSize,result)
		}
	}else{
		if curId!=""{
			return rev_prevPage(c,query,pageSize,curId,result)
		}else{
			return rev_endPage(c,query,pageSize,result)
		}
	}
}

func RevPage(c *mgo.Collection,query map[string]interface{},direction string,curId string,isCount bool,pageSize int,result interface{})(*PageBean,error){
	if pageSize<=0{
		pageSize=10
	}
	total:=0
	var err error
	if isCount==true{
		total,err=c.Find(query).Count();if err!=nil{
			return nil,err
		}
	}


	//开始查询
	pageBean,err:=rev_queryPageBean(c,query,direction,pageSize,curId,result);if err!=nil{
		return nil,err
	}
	pageBean.Total=total
	if total>0{
		//分页
		pageBean.TotalPage=int(math.Ceil(float64(total)/float64(pageSize)))
	}

	return pageBean,nil
}
