package mgowrap

import (
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PersistentObject interface {
	CollectionName() string
	MakeId() interface{}
	SetTimeStamp()
}

func (db *Database) Save(po PersistentObject) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		po.SetTimeStamp()
		_, err = rc.Upsert(bson.M{"_id": po.MakeId()}, po)
	})
	return
}

func (db *Database) SaveAll(items []interface{}) (err error) {
	if len(items) == 0 {
		return
	}
	item := reflect.ValueOf(items[0]).Elem()
	db.CollectionDo(callCollectionName(item), func(rc *mgo.Collection) {
		err = rc.Insert(items...)
	})

	return
}

func (db *Database) FindById(id bson.ObjectId, result interface{}) (err error) {
	return db.Find(bson.M{"_id": id}, result)
}

func (db *Database) FindByIdHex(idHex string, result interface{}) (err error) {
	id, err := ToObjectId(idHex)
	if err != nil {
		return
	}
	return db.FindById(id, result)
}

func (db *Database) Find(query, result interface{}) (err error) {
	item := reflect.ValueOf(result).Elem()
	db.CollectionDo(callCollectionName(item), func(c *mgo.Collection) {
		err = c.Find(query).One(result)
	})
	return
}

func (db *Database) FindAll(query interface{}, result interface{}, sortFields ...string) (err error) {
	resultv := reflect.ValueOf(result)
	resultvKind := resultv.Kind()
	slicev := resultv.Elem()

	if resultvKind != reflect.Ptr || slicev.Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	element := slicev.Type().Elem().Elem()
	newValue := reflect.New(element)

	db.CollectionDo(callCollectionName(newValue), func(c *mgo.Collection) {
		if len(sortFields) == 0 {
			err = c.Find(query).All(result)
		} else {
			err = c.Find(query).Sort(sortFields...).All(result)
		}
	})
	return
}

func (db *Database) Update(po PersistentObject, selector, changer interface{}) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Update(selector, changer)
	})
	return
}

func (db *Database) UpdateInstance(po PersistentObject, changer interface{}) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Update(bson.M{"_id": po.MakeId()}, changer)
	})
	return
}

func (db *Database) UpdateAll(po PersistentObject, selector, changer interface{}) (info *mgo.ChangeInfo, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		info, err = rc.UpdateAll(selector, changer)
	})
	return
}

func (db *Database) Delete(po PersistentObject, selector interface{}) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Remove(selector)
	})
	return
}

func (db *Database) DeleteAll(po PersistentObject, selector interface{}) (info *mgo.ChangeInfo, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		info, err = rc.RemoveAll(selector)
	})
	return
}

func (db *Database) Count(po PersistentObject, selector bson.M) (count int, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		count, err = rc.Find(selector).Count()
	})
	return
}

func callCollectionName(value reflect.Value) string {
	method := value.MethodByName("CollectionName")
	if !method.IsValid() {
		panic(value.String() + ` does not implement "CollectionName" method`)
	}

	return method.Call([]reflect.Value{})[0].String()
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