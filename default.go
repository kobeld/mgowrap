package mgowrap

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var defaultDatabase *Database

func DefaultDatabase() *Database {
	return defaultDatabase
}

func SetupDatbase(dialString string, name string) {
	defaultDatabase = NewDatabase(dialString, name)
}

func DropDatabase() (err error) {
	return defaultDatabase.DropDatabase()
}

func DropCollection(name string) (err error) {
	defaultDatabase.CollectionDo(name, func(rc *mgo.Collection) {
		err = rc.DropCollection()
	})
	return
}

func DatabaseDo(f func(db *mgo.Database)) {
	defaultDatabase.DatabaseDo(f)
}

func CollectionDo(name string, f func(c *mgo.Collection)) {
	defaultDatabase.CollectionDo(name, f)
}

func CollectionsDo(f func(c ...*mgo.Collection), names ...string) {
	defaultDatabase.CollectionsDo(f, names...)
}

// Curd functions

func Save(po PersistentObject, funcs ...func()) error {
	for _, f := range funcs {
		f()
	}
	return defaultDatabase.Save(po)
}

func SaveAll(items []interface{}) error {
	return defaultDatabase.SaveAll(items)
}

func FindById(id bson.ObjectId, result interface{}) (err error) {
	return defaultDatabase.Find(bson.M{"_id": id}, result)
}

func FindByIdHex(idHex string, result interface{}) (err error) {
	id, err := ToObjectId(idHex)
	if err != nil {
		return
	}
	return defaultDatabase.FindById(id, result)
}

func Find(query, result interface{}) error {
	return defaultDatabase.Find(query, result)
}

func FindAll(query interface{}, result interface{}, sortFilelds ...string) error {
	return defaultDatabase.FindAll(query, result, sortFilelds...)
}

func Upsert(po PersistentObject, selector, changer interface{}) (*mgo.ChangeInfo, error) {
	return defaultDatabase.Upsert(po, selector, changer)
}

func Update(po PersistentObject, selector, changer interface{}) error {
	return defaultDatabase.Update(po, selector, changer)
}

func UpdateInstance(po PersistentObject, changer interface{}) error {
	return defaultDatabase.UpdateInstance(po, changer)
}

func UpdateAll(po PersistentObject, selector, changer interface{}) (*mgo.ChangeInfo, error) {
	return defaultDatabase.UpdateAll(po, selector, changer)
}

func Delete(po PersistentObject, selector interface{}) error {
	return defaultDatabase.Delete(po, selector)
}

func DeleteAll(po PersistentObject, selector interface{}) (*mgo.ChangeInfo, error) {
	return defaultDatabase.DeleteAll(po, selector)
}

func Count(po PersistentObject, selector bson.M) (int, error) {
	return defaultDatabase.Count(po, selector)
}
