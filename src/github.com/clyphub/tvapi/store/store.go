package store

import (
	"errors"
	"fmt"
	"github.com/clyphub/tvapi/objects"
	"reflect"
	"regexp"
)

type ObjectStore interface {
	Save(obj objects.Storable) (string, error)
	Get(key string) (objects.Storable, error)
	Erase(key string) (objects.Storable, error)
	Find(fieldName string, fieldVal string) ([]objects.Storable, error)
}

type MapStore struct {
	store map[string]objects.Storable
}

func NewMapStore() *MapStore {
	return &MapStore{store: make(map[string]objects.Storable,10)}
}

func(s *MapStore) Save(obj objects.Storable) (string, error) {
	if(obj == nil) {
		return "", errors.New("Object is nil")
	}
	k := obj.GetKey()
	if(len(k) == 0) {
		return "", errors.New("Key is nil")
	}

	s.store[k] = obj
	return k, nil
}

func(s *MapStore) Get(key string) (objects.Storable, error){
	if(len(key) == 0) {
		return nil, errors.New("Key is nil")
	}
	return s.store[key], nil
}

func(s *MapStore) Erase(key string) (objects.Storable, error){
	if(len(key)==0) {
		return nil, errors.New("Key is nil")
	}
	o := s.store[key]
	delete(s.store, key)
	return o, nil
}

func(s *MapStore) Find(fieldName string, fieldValPattern string) ([]objects.Storable, error) {
	if(len(fieldName) == 0) {
		return nil, errors.New("Key is nil")
	}
	if(len(fieldValPattern) == 0) {
		return nil, nil
	}
	regex, rexerr := regexp.Compile(fieldValPattern)
	if(rexerr != nil){
		return nil, rexerr
	}

	res := make([]objects.Storable, 0)
	for _, val := range s.store {
		valStr, e := s.getField(fieldName, val)
		if(e != nil){
			return nil, e
		}
		if(regex.MatchString(valStr)){
			res = append(res, val)
		}
	}
	return res, nil
}

func(s *MapStore) getField(fieldName string, obj objects.Storable) (string, error) {
	val := reflect.ValueOf(obj).Elem()
	f := val.FieldByName(fieldName)
	return fmt.Sprintf("%v", f.Interface()), nil
}

