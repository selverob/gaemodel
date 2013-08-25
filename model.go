package gaemodel

import (
	"appengine/datastore"
)

type Model interface {
	Kind() string

	Key() *datastore.Key
	SetKey(*datastore.Key)

	Ancestor() *datastore.Key
}
