package mgowrap

import (
	"errors"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PersistentObject interface {
	CollectionName() string
	MakeId() interface{}
}

func (db *Database) EnsureIndexKey(po PersistentObject, keys ...string) (err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		err = rc.EnsureIndexKey(keys...)
	})
	return
}

func (db *Database) Save(po PersistentObject, funcs ...func()) (err error) {

	for _, f := range funcs {
		f()
	}

	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {

		_, err = rc.Upsert(bson.M{"_id": po.MakeId()}, po)
	})
	return
}

func (db *Database) SaveAll(items []interface{}) (err error) {
	if len(items) == 0 {
		return
	}
	item := reflect.ValueOf(items[0])
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

	name, err := callCollectionNameForSingleItem(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {
		err = c.Find(query).One(result)
	})
	return
}

func (db *Database) FindAndSelect(query, selector, result interface{}) (err error) {

	name, err := callCollectionNameForSingleItem(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {
		err = c.Find(query).Select(selector).One(result)
	})
	return
}

func (db *Database) FindAllAndSelect(query, selector, result interface{}) (err error) {

	name, err := callCollectionNameForItems(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {
		err = c.Find(query).Select(selector).All(result)
	})
	return
}

func (db *Database) FindAndApply(query interface{}, change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	name, err := callCollectionNameForSingleItem(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {
		info, err = c.Find(query).Apply(change, result)
	})
	return
}

func (db *Database) FindAll(query interface{}, result interface{}, sortFields ...string) (err error) {
	name, err := callCollectionNameForItems(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {
		if len(sortFields) == 0 {
			err = c.Find(query).All(result)
		} else {
			err = c.Find(query).Sort(sortFields...).All(result)
		}
	})
	return
}

func (db *Database) FindWithLimit(selector interface{}, result interface{}, limit int, sortFields ...string) (err error) {
	name, err := callCollectionNameForItems(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {

		query := c.Find(selector)

		// Ensure the sort fields are not empty
		var validFields = []string{}
		for _, field := range sortFields {
			if str := strings.TrimSpace(field); str != "" {
				validFields = append(validFields, strings.ToLower(str))
			}
		}
		if len(validFields) > 0 {
			query.Sort(validFields...)
		}

		if limit > 0 {
			query.Limit(limit)
		}

		err = query.All(result)
	})
	return
}

func (db *Database) FindWithSkipAndLimit(selector interface{}, result interface{}, skip, limit int, sortFields ...string) (err error) {
	name, err := callCollectionNameForItems(result)
	if err != nil {
		return
	}

	db.CollectionDo(name, func(c *mgo.Collection) {

		query := c.Find(selector)

		// Ensure the sort fields are not empty
		var validFields = []string{}
		for _, field := range sortFields {
			if str := strings.TrimSpace(field); str != "" {
				validFields = append(validFields, strings.ToLower(str))
			}
		}
		if len(validFields) > 0 {
			query.Sort(validFields...)
		}

		if skip > 0 {
			query.Skip(skip)
		}

		if limit > 0 {
			query.Limit(limit)
		}

		err = query.All(result)
	})
	return
}

func (db *Database) Upsert(po PersistentObject, selector, changer interface{}) (changeInfo *mgo.ChangeInfo, err error) {
	db.CollectionDo(po.CollectionName(), func(rc *mgo.Collection) {
		changeInfo, err = rc.Upsert(selector, changer)
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

func (db *Database) DeleteInstance(po PersistentObject) (err error) {
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

func callCollectionNameForSingleItem(result interface{}) (name string, err error) {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		err = errors.New("Result argument must be a pointer")
		return
	}

	item := resultv.Elem()
	if item.Kind() != reflect.Ptr && item.Kind() != reflect.Struct {
		err = errors.New("Result argument must point to a struct or struct pointer")
		return
	}

	name = callCollectionName(item)

	return
}

func (db *Database) DropCollection(collectionName string) (err error) {
	db.CollectionDo(collectionName, func(rc *mgo.Collection) {
		err = rc.DropCollection()
	})
	return
}

// ======
// ====== Private =====
// ======
func callCollectionNameForItems(result interface{}) (name string, err error) {
	resultv := reflect.ValueOf(result)

	if resultv.Kind() != reflect.Ptr {
		err = errors.New("Result argument must be a pointer to slice")
		return
	}

	slicev := resultv.Elem()
	if slicev.Kind() != reflect.Slice {
		err = errors.New("Result argument must pointer to a slice")
		return
	}

	element := slicev.Type().Elem().Elem()
	newValue := reflect.New(element)

	name = callCollectionName(newValue)

	return
}

func callCollectionName(value reflect.Value) string {
	method := value.MethodByName("CollectionName")
	if !method.IsValid() {
		panic(value.String() + ` does not implement "CollectionName" method`)
	}

	return method.Call([]reflect.Value{})[0].String()
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
