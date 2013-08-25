package gaemodel

import (
	"appengine"
	"appengine/datastore"
	"math"
	"reflect"
)

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

func Delete(c appengine.Context, m Model) (err error) {
	err = datastore.Delete(c, m.Key())
	if err != nil {
		return err
	}
	m.SetKey(nil)
	return
}

func PageCount(c appengine.Context, kind string, perPage int) (pages int, err error) {
	count, err := datastore.NewQuery(kind).Count(c)
	if err != nil {
		return
	}
	pages = int(math.Ceil(float64(count) / float64(perPage)))
	return
}

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

func GetByAncestor(c appengine.Context, typ reflect.Type, kind string, anc *datastore.Key) (interface{}, error) {
	query := datastore.NewQuery(kind).Ancestor(anc)
	return MultiQuery(c, typ, kind, query)
}

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
