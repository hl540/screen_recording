package model

import (
	"fmt"
)

type MessageEvent struct {
	Id    string
	Event string
	Data  string
}

func (e MessageEvent) String() string {
	f := "id:%s\n"
	f += "event:%s\n"
	f += "data:%s\n\n"
	return fmt.Sprintf(f, e.Id, e.Event, e.Data)
}
