package gaemodel

import (
	"appengine"
	"appengine/datastore"
	"math"
	"reflect"
)

//Save saves a Model into Datastore. If its key is nil, it inserts it.
//If its key is set, it replaces the entity with that key with given Model.
func Save(c appengine.Context, m Model) (err error) {
	if m.Key() == nil {
		var k *datastore.Key
		k, err = datastore.Put(c, datastore.NewIncompleteKey(c, m.Kind(), m.Ancestor()), m)
		if err != nil {
			return
		}
		m.SetKey(k)
	} else {
		_, err = datastore.Put(c, m.Key(), m)
	}
	return
}

//Delete removes given Model from Datastore. Since it only
//compares by keys, you don't have to load the whole entity
//if you only want to delete it. You only need to set the key
//of an empty struct:
//	key := GetKeyOfEmployeeToDelete()
//	e := new(Employee)
//	e.SetKey(key)
//	err := e.Delete(c)
func Delete(c appengine.Context, m Model) (err error) {
	err = datastore.Delete(c, m.Key())
	if err != nil {
		return err
	}
	m.SetKey(nil)
	return
}

//PageCount returns number of pages which would be needed when listing
//all of the entities of given kind, listing perPage entries
//per page.
func PageCount(c appengine.Context, kind string, perPage int) (pages int, err error) {
	count, err := datastore.NewQuery(kind).Count(c)
	if err != nil {
		return
	}
	pages = int(math.Ceil(float64(count) / float64(perPage)))
	return
}

//GetByKey returns entity from Datastore with the given key.
//It wraps datastore.Get function.
func GetByKey(c appengine.Context, typ reflect.Type, key *datastore.Key) (m interface{}, err error) {
	v := reflect.New(typ)
	err = datastore.Get(c, key, v.Interface())
	if err != nil {
		return
	}
	v.MethodByName("SetKey").Call([]reflect.Value{reflect.ValueOf(key)})
	m = v.Interface()
	return
}

//GetAll returns a paged listing of all the entities of given kind.
//If page is 0, perPage is ignored and it returns all the entities of the given kind.
func GetAll(c appengine.Context, typ reflect.Type, kind string, page, perPage int) (interface{}, error) {
	var offset, limit int
	if page == 0 {
		limit = -1
		offset = 0
	} else {
		offset = (page - 1) * perPage
		limit = perPage
	}

	query := datastore.NewQuery(kind).Limit(limit).Offset(offset)

	return MultiQuery(c, typ, kind, query)
}

//GetByAncestor returns all the entities with given ancestor.
func GetByAncestor(c appengine.Context, typ reflect.Type, kind string, anc *datastore.Key) (interface{}, error) {
	query := datastore.NewQuery(kind).Ancestor(anc)
	return MultiQuery(c, typ, kind, query)
}

//MultiQuery executes given query and returns slice of all the entities it returns, with their keys set.
func MultiQuery(c appengine.Context, typ reflect.Type, kind string, query *datastore.Query) (ms interface{}, err error) {
	is := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(typ)), 0, 0)

	iter := query.Run(c)

	for {
		val := reflect.New(typ)
		var key *datastore.Key
		key, err = iter.Next(val.Interface())
		if err != nil {
			if err == datastore.Done {
				err = nil
				val.MethodByName("SetKey").Call([]reflect.Value{reflect.ValueOf(key)})
				reflect.Append(is, val)
				break
			}
			return
		}
		val.MethodByName("SetKey").Call([]reflect.Value{reflect.ValueOf(key)})
		is = reflect.Append(is, val)
	}

	ms = is.Interface()

	return
}
