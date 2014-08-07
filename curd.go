package mgowrap

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type PersistentObject interface {
	CollectionName() string
	MakeId() interface{}
}

func (db *Database) Save(po PersistentObject) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		_, err = rc.Upsert(bson.M{"_id": po.MakeId()}, po)
	})
	return
}

func (db *Database) FindAll(collectionName string, query interface{}, result interface{}) (err error) {
	db.CollectionDo(collectionName, func(c *mgo.Collection) {
		err = c.Find(query).All(result)
	})
	return
}

// func (db *Database) DropCollection(collectionName string) (err error) {
// 	db.CollectionDo(collectionName, func(rc *mgo.Collection) {
// 		err = rc.DropCollection()
// 	})
// 	return
// }

// func (db *Database) DropCollections(collectionNames ...string) (err error) {
// 	db.CollectionsDo(func(rcs ...*mgo.Collection) {
// 		for _, rc := range rcs {
// 			err1 := rc.DropCollection()
// 			if err == nil && err1 != nil {
// 				err = err1
// 			}
// 		}
// 	}, collectionNames...)
// 	return
// }

// func (db *Database) Delete(collectionName string, id interface{}) (err error) {
// 	db.CollectionDo(collectionName, func(rc *mgo.Collection) {
// 		_, err = rc.RemoveAll(bson.M{"_id": id})
// 	})
// 	return
// }

// func (db *Database) Update(collectionName string, obj Id) (err error) {
// 	db.CollectionDo(collectionName, func(rc *mgo.Collection) {
// 		v := reflect.ValueOf(obj)
// 		if v.Kind() == reflect.Ptr {
// 			v = v.Elem()
// 		}

// 		found := reflect.New(v.Type()).Interface()
// 		rc.Find(bson.M{"_id": obj.MakeId()}).One(found)

// 		originalValue := reflect.ValueOf(found)
// 		if originalValue.Kind() == reflect.Ptr {
// 			originalValue = originalValue.Elem()
// 		}

// 		for i := 0; i < v.NumField(); i++ {
// 			fieldValue := v.Field(i)
// 			if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
// 				continue
// 			}

// 			fieldValue.Set(originalValue.Field(i))
// 		}

// 		rc.Upsert(bson.M{"_id": obj.MakeId()}, v.Interface())
// 	})
// 	return
// }

// func (db *Database) FindOne(collectionName string, query interface{}, result interface{}) (err error) {
// 	db.CollectionDo(collectionName, func(c *mgo.Collection) {
// 		err = c.Find(query).One(result)
// 	})
// 	return
// }

// func (db *Database) FindById(collectionName string, id interface{}, result interface{}) (err error) {
// 	db.CollectionDo(collectionName, func(c *mgo.Collection) {
// 		err = c.Find(bson.M{"_id": id}).One(result)
// 	})
// 	return
// }
