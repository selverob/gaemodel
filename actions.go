package gaemodel

import (
	"appengine"
	"appengine/datastore"
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

func GetAll(c appengine.Context, kind string, page, perPage int) (ms []Model, err error) {
	var offset int
	if page == 0 {
		perPage = -1
		offset = 0
	} else {
		offset = (page - 1) * perPage
	}

	ms = make([]Model, 0)
	keys, err := datastore.NewQuery(kind).Limit(perPage).Offset(offset).GetAll(c, &ms)
	if err != nil {
		return
	}
	for i, k := range keys {
		ms[i].SetKey(k)
	}
	return
}
