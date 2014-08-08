package mgowrap

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var defaultDatabase *Database

func SetupDatbase(dialString string, name string) {
	defaultDatabase = NewDatabase(dialString, name)
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

func Find(query, result interface{}) error {
	return defaultDatabase.Find(query, result)
}

func FindAll(query interface{}, result interface{}, sortFilelds ...string) error {
	return defaultDatabase.FindAll(query, result, sortFilelds...)
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
