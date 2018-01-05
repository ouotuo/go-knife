package redis_utils

import (
	"net/url"
	"gopkg.in/redis.v5"
)

func Connect(redisUrl string)(client *redis.Client,err error){
	u,err:=url.Parse(redisUrl)
	if err!=nil{
		return
	}

	p:=""
	if u.User!=nil{
		p,_=u.User.Password()
	}
	client = redis.NewClient(&redis.Options{
		Addr:    u.Host,
		Password: p,
		DB:       0,  // use default DB
	})
	sc:=client.Ping()
	if sc.Err()!=nil{
		err=sc.Err()
	}

	return
}
