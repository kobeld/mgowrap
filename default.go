package mgowrap

import (
	"fmt"
	"labix.org/v2/mgo"
)

var defaultDatabase *Database

func SetupDatbase(dialString string, name string) {
	defaultDatabase = NewDatabase(dialString, name)
	fmt.Printf("******** %+v \n", defaultDatabase)
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

func Save(po PersistentObject) error {
	return defaultDatabase.Save(po)
}
