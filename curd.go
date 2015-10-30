package mgowrap

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	if err == mgo.ErrNotFound {
		err = nil
	}

	return
}

func (db *Database) FindAndSelect(query, selector, result interface{}) (err error) {
	item := reflect.ValueOf(result).Elem()
	db.CollectionDo(callCollectionName(item), func(c *mgo.Collection) {
		err = c.Find(query).Select(selector).One(result)
	})

	if err == mgo.ErrNotFound {
		err = nil
	}

	return
}

func (db *Database) FindAndApply(query interface{}, change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	item := reflect.ValueOf(result).Elem()
	db.CollectionDo(callCollectionName(item), func(c *mgo.Collection) {
		info, err = c.Find(query).Apply(change, result)
	})

	if err == mgo.ErrNotFound {
		err = nil
	}
	return
}

func (db *Database) FindAll(query interface{}, result interface{}, sortFields ...string) (err error) {
	resultv := reflect.ValueOf(result)
	resultvKind := resultv.Kind()

	if resultvKind != reflect.Ptr {
		return errors.New("Result argument must be a pointer to slice")
	}

	slicev := resultv.Elem()
	if slicev.Kind() != reflect.Slice {
		return errors.New("Result argument must be a pointer to slice")
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

func (db *Database) Upsert(po PersistentObject, selector, changer interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		changeInfo, err = rc.Upsert(selector, changer)
	})
	return
}

func (db *Database) Update(po PersistentObject, selector, changer interface{}) (ok bool, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Update(selector, changer)
	})

	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *Database) UpdateInstance(po PersistentObject, changer interface{}) (ok bool, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Update(bson.M{"_id": po.MakeId()}, changer)
	})

	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *Database) UpdateAll(po PersistentObject, selector, changer interface{}) (info *mgo.ChangeInfo, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		info, err = rc.UpdateAll(selector, changer)
	})
	return
}

func (db *Database) Delete(po PersistentObject, selector interface{}) (ok bool, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.Remove(selector)
	})

	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *Database) DeleteInstance(po PersistentObject) (ok bool, err error) {
	return db.Delete(po, bson.M{"_id": po.MakeId()})
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

func (db *Database) HasAny(po PersistentObject, selector bson.M) (r bool, err error) {
	count, err := db.Count(po, selector)
	if err != nil {
		return
	}

	r = (count > 0)
	return
}

func callCollectionName(value reflect.Value) string {
	method := value.MethodByName("CollectionName")
	if !method.IsValid() {
		panic(value.String() + ` does not implement "CollectionName" method`)
	}

	return method.Call([]reflect.Value{})[0].String()
}

func (db *Database) DropCollection(collectionName string) (err error) {
	db.CollectionDo(collectionName, func(rc *mgo.Collection) {
		err = rc.DropCollection()
	})
	return
}

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
