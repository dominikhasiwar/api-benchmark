package entities

import (
	"time"
)

type IEntityBase interface {
	GetId() string
	SetId(id string)
	GetCreator() string
	SetCreator(creator string)
	GetCreated() time.Time
	SetCreated(created time.Time)
	GetModifier() string
	SetModifier(modifier string)
	GetModified() time.Time
	SetModified(modified time.Time)
}

type EntityBase struct {
	Id       string    `dynamodbav:"Id"`
	Creator  string    `dynamodbav:"Creator"`
	Created  time.Time `dynamodbav:"Created"`
	Modifier string    `dynamodbav:"Modifier"`
	Modified time.Time `dynamodbav:"Modified"`
}

func (e *EntityBase) GetId() string {
	return e.Id
}

func (e *EntityBase) SetId(id string) {
	e.Id = id
}

func (e *EntityBase) GetCreator() string {
	return e.Creator
}

func (e *EntityBase) SetCreator(creator string) {
	e.Creator = creator
}

func (e *EntityBase) GetCreated() time.Time {
	return e.Created
}

func (e *EntityBase) SetCreated(created time.Time) {
	e.Created = created
}

func (e *EntityBase) GetModifier() string {
	return e.Modifier
}

func (e *EntityBase) SetModifier(modifier string) {
	e.Modifier = modifier
}

func (e *EntityBase) GetModified() time.Time {
	return e.Modified
}

func (e *EntityBase) SetModified(modified time.Time) {
	e.Modified = modified
}
