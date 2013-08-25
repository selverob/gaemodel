//Gaemodel is a library that simplifies working with App Engine's Datastore.
//It only needs a few modifications to your models:
//If your model is a plain struct, like the one from Datastore introduction:
//	type Employee struct {
//		Name     string
//		Role     string
//		HireDate time.Time
//		Account  string
//	}
//You need to add a single field to it:
//	type Employee struct {
//		Name     string
//		Role     string
//		HireDate time.Time
//		Account  string
//		key      *datastore.Key `datastore:"-"`
//	}
//And implement methods from Model interface:
//	func (e *Employee) Key() *datastore.Key {
//		return 	e.key
//	}
//
//	func (e *Employee) SetKey(k *datastore.Key) {
//		e.key = k
//	}
//	//and so on...
//After this, you can wrap functions provided by gaemodel as struct's methods:
//	func (e *Employee) Save(c appengine.Context) (err error) {
//		return gaemodel.Save(c, e)
//	}
//You have to be careful, though, to correctly convert results of functions:
//	//These are query functions, not methods
//	func GetByKey(c appengine.Context, k *datastore.Key) (e *Employee, err error) {
//		m, err := gaemodel.GetByKey(c, typ, k)
//		if err != nil {
//			return
//		}
//		o = m.(*Employee)
//		return
//	}
//
//	func GetAll(c appengine.Context, page, perPage int) (es []*Employee, err error) {
//		ms, err := gaemodel.GetAll(c, typ, "Employee", page, perPage)
//		if err != nil {
//			return
//		}
//		es = ms.([]*Employee)
//		return
//	}
package gaemodel

import (
	"appengine/datastore"
)

//Model is the interface structs used by gaemodel must implement.
//Be careful, because even the query functions that don't have Model
//mentioned in their signatures only support querying types that implement it.
type Model interface {
	Kind() string

	Key() *datastore.Key
	SetKey(*datastore.Key)

	Ancestor() *datastore.Key
}
