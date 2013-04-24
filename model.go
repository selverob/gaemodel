package gaemodel

import (
	"appengine/datastore"
)

type Model interface {
	GetKind() string

	GetKey() *datastore.Key
	SetKey(*datastore.Key)

	GetAncestor() *datastore.Key
}
