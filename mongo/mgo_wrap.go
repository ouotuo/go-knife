package mongo

import (
	"gopkg.in/mgo.v2"
	"time"
	"log"
)

type MgoWrap struct{
	session *mgo.Session
}
func(self *MgoWrap)SetDebug(){
	mgo.SetDebug(true)
	mgo.SetLogger(self)
}

func(self *MgoWrap)Output(calldepth int, s string) error{
	log.Println(s)
	return nil
}



func NewMgoWrap(url string)(*MgoWrap,error){
	self:=&MgoWrap{}

	var err error
	self.session, err = mgo.DialWithTimeout(url,time.Second*5)
	if err != nil {
		return nil,err
	}

	//self.session.SetMode(mgo.Monotonic,true)
	self.session.SetMode(mgo.Eventual,true)

	return self,nil
}

func(self *MgoWrap)SetMode(consistency mgo.Mode, refresh bool){
	self.session.SetMode(consistency,refresh);
}

func (self *MgoWrap)NewSession()(*mgo.Session){
	session:=self.session.Copy()

	return session
}

func (self *MgoWrap)Close(){
	self.session.Close()
}