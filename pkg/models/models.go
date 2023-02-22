package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Username  string             `json:"username,omitempty" bson:"username"`
	Password  string             `json:"password" bson:"password"`
}

type Item struct {
	Id          string                 `json:"_id" bson:"-"`
	PrimitiveId primitive.ObjectID     `json:"-" bson:"_id,omitempty"`
	Owner       string                 `json:"username,omitempty" bson:"owner,omitempty"`
	Created     time.Time              `json:"created" bson:"created,omitempty"`
	Modified    time.Time              `json:"modified" bson:"modified"`
	Data        map[string]interface{} `json:"data" bson:"data,omiempty"`
}

type UpdateItem struct {
	Id       string                 `bson:"-"`
	Modified time.Time              `bson:"modifed"`
	Data     map[string]interface{} `bson:"data"`
}

func (u *UpdateItem) Update(i Item) bson.M {
	update := bson.M{}

	for k,v := range i.Data{
		_,ok := u.Data[k]
		if !ok{
			u.Data[k] =v
		}
	}
	update["modified"] = time.Now()
	update["data"] = u.Data
	return bson.M{"$set": update}
}


