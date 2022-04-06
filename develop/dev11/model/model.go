package model

import (
	//"fmt"
	//"strings"
	"time"
)

//Event struct
type Event struct {
	//EventID int    `json:"event_id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Descr   string `json:"descr"`
	Date    Date   `json:"date"`
}

type Date struct {
	time.Time
}

