package mongo

import (
	"time"
	"fmt"
)

type JsonTime time.Time

func (t JsonTime)MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}


