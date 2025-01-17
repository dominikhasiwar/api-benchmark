package models

import "time"

type BaseModel struct {
	Id       string    `json:"id"`
	Creator  string    `json:"creator"`
	Created  time.Time `json:"created"`
	Modifier string    `json:"modifier"`
	Modified time.Time `json:"modified"`
}
