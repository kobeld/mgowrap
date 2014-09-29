package mgowrap

import (
	"testing"
)

const (
	DIAL_STRING = "localhost"
	DB_NAME     = "mgowrap_testing"
)

type User struct {
	Id    bson.ObjectId `bson:"_id"`
	Email string
}

func (this *User) MakeId() interface{} {
	if this.Id == "" {
		this.Id = bson.NewObjectId()
	}
	return this.Id
}

func (this *User) CollectionName() string {
	return "users"
}

func TestSave(t *testing.T) {

	db := NewDatabase(DIAL_STRING, DB_NAME)

	aaron := &User{
		Email: "aaron@theplant.jp",
	}

	err := db.Save(aaron)
	if err != nil {
		t.Error(err)
	}

	// db.Save(ALLUSERS, &User{Email: "sunfmin@gmail.com", Name: "Felix Sun"})

	// var found *User
	// db.CollectionDo(ALLUSERS, func(uc *mgo.Collection) {
	// 	uc.Find(bson.M{"email": "sunfmin@gmail.com"}).One(&found)
	// })
	// if found == nil {
	// 	t.Error("Can not find user after saved")
	// }

	// db = NewDatabase("localhost", DB2)

	// db.Save(ALLUSERS, &User{Email: "sunfmin@gmail.com", Name: "Felix Sun"})

	// var u2 *User
	// db.CollectionDo(ALLUSERS, func(uc *mgo.Collection) {
	// 	uc.Find(bson.M{"email": "sunfmin@gmail.com"}).One(&u2)
	// })
	// if u2 == nil {
	// 	t.Error("Can not find user after saved")
	// }

}
