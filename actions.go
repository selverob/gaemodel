package gaemodel

import (
	"appengine"
	"appengine/datastore"
	"math"
)

func Save(c appengine.Context, m Model) (err error) {
	if m.GetKey() == nil {
		var k *datastore.Key
		k, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Owner", m.GetAncestor()), m)
		if err != nil {
			return
		}
		m.SetKey(k)
	} else {
		_, err = datastore.Put(c, m.GetKey(), m)
	}
	return
}

func Delete(c appengine.Context, m Model) (err error) {
	err = datastore.Delete(c, m.GetKey())
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
