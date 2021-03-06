package mgowrap

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var defaultDatabase *Database

func DefaultDatabase() *Database {
	return defaultDatabase
}

func SetupDatabase(dialString string, name string) {
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

func EnsureIndexKey(po PersistentObject, keys ...string) {
	defaultDatabase.EnsureIndexKey(po, keys...)
}

// Curd functions

func Save(po PersistentObject, funcs ...func()) error {
	return defaultDatabase.Save(po, funcs...)
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

func FindAndSelect(query, selector, result interface{}) error {
	return defaultDatabase.FindAndSelect(query, selector, result)
}

func FindAndApply(query interface{}, change mgo.Change, result interface{}) (*mgo.ChangeInfo, error) {
	return defaultDatabase.FindAndApply(query, change, result)
}

func FindAll(query interface{}, result interface{}, sortFilelds ...string) error {
	return defaultDatabase.FindAll(query, result, sortFilelds...)
}

func FindAllAndSelect(query, selector, result interface{}, sortFields ...string) error {
	return defaultDatabase.FindAllAndSelect(query, selector, result, sortFields...)
}

func FindWithLimit(query interface{}, result interface{}, limit int, sortFilelds ...string) error {
	return defaultDatabase.FindWithLimit(query, result, limit, sortFilelds...)
}

func FindWithSkipAndLimit(query interface{}, result interface{}, skip, limit int, sortFilelds ...string) error {
	return defaultDatabase.FindWithSkipAndLimit(query, result, skip, limit, sortFilelds...)

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

func DeleteInstance(po PersistentObject) error {
	return defaultDatabase.DeleteInstance(po)
}

func DeleteAll(po PersistentObject, selector interface{}) (*mgo.ChangeInfo, error) {
	return defaultDatabase.DeleteAll(po, selector)
}

func Count(po PersistentObject, selector bson.M) (int, error) {
	return defaultDatabase.Count(po, selector)
}

func HasAny(po PersistentObject, selector bson.M) (bool, error) {
	return defaultDatabase.HasAny(po, selector)
}
