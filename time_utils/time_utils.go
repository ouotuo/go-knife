package time_utils

import (
	"time"
	"fmt"
)


const (
	FORMAT_LONG="2006-01-02 15:04:05"
	FORMAT_SHORT="2006-01-02"
	FORMAT_SHORT1="20060102"

	CONVERT_MS_NS=1000000
	SECOND_MS=1000
	MINUTE_MS=SECOND_MS*60
	HOUR_MS=MINUTE_MS*60
	DAY_MS=HOUR_MS*24
)



func GetNowMs()int64{
	return time.Now().UnixNano()/CONVERT_MS_NS
}

func GetNowLong()string{
	return time.Now().Format(FORMAT_LONG)
}

func GetNowShort()string{
	return time.Now().Format(FORMAT_SHORT)
}

func GetNowMsAndLong()(int64,string){
	now:=time.Now()

	return now.UnixNano()/CONVERT_MS_NS,now.Format(FORMAT_LONG)
}

func FormatShort(t time.Time)(string){
	return t.Format(FORMAT_SHORT)
}

func GetEndOfDay(t time.Time)time.Time{
	shortStr:=FormatShort(t)
	longStr:=fmt.Sprintf("%s 23:59:59",shortStr)
	rt,_:=time.ParseInLocation(FORMAT_LONG,longStr,time.Local)
	return rt
}

func GetTimeOfDay(t time.Time,hms string)(time.Time,error){
	shortStr:=FormatShort(t)
	longStr:=fmt.Sprintf("%s %s",shortStr,hms)
	return time.ParseInLocation(FORMAT_LONG,longStr,time.Local)
}

//date 20170118 time 150618  转为time
func ParseFromDateAndTime(d int,t int)(rt time.Time,err error){
	rt,err=time.ParseInLocation(FORMAT_LONG,fmt.Sprintf("%d-%.2d-%.2d %.2d:%.2d:%.2d",d/10000,d/100%100,d%100,t/10000,t/100%100,t%100),time.Local)
	return
}

//从长格式转换为time
func ParseFromLongStr(t string)(rt time.Time,err error){
	rt,err=time.ParseInLocation(FORMAT_LONG,t,time.Local)
	return
}

func ParseFromShortStr(t string)(rt time.Time,err error){
	rt,err=time.ParseInLocation(FORMAT_SHORT,t,time.Local)
	return
}

//从ms转换为time
func ParseFromMs(t int64)(rt time.Time,err error){
	if t<=0{
		err=fmt.Errorf("time is <=0")
		return
	}
	rt=time.Unix(t/1000,0)
	return
}

func IsSameDay(t1 time.Time,t2 time.Time)(bool){
	return FormatShort(t1)==FormatShort(t2)
}

func GetBeginTimeOfDayNowMs()(bt int64){
	now:=time.Now()
	t,err:=GetTimeOfDay(now,"00:00:00")
	if err!=nil{
		t=now
	}

	return t.UnixNano()/CONVERT_MS_NS
}

func ParseFromMsNoError(t int64)(rt time.Time){
	rt=time.Unix(t/1000,0)
	return
}
