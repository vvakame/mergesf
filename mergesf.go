package mergesf

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var rootCache = newReflectCache()

var RecoverPanic bool

func Merge(objs ...interface{}) (_ interface{}, err error) {
	if v := len(objs); v == 0 {
		return nil, nil
	} else if v == 1 {
		return objs[0], nil
	}

	if RecoverPanic {
		defer func() {
			rErr := recover()
			if err != nil {
				return
			}
			if rErr == nil {
				return
			}
			err = fmt.Errorf("%s", rErr)
		}()
	}

	ktc, err := rootCache.getTypeCache(objs, nil)
	if err != nil {
		return nil, err
	}

	return ktc.mergeObjects(objs)
}

type reflectCache struct {
	sync.RWMutex
	nextTypes map[reflect.Type]*reflectCache
	typeCache *knownTypeCache
}

type knownTypeCache struct {
	st  reflect.Type
	fas []*fieldAccessor
}

type fieldAccessor struct {
	toFieldIndex int

	objIndex   int
	fieldIndex int
}

func newReflectCache() *reflectCache {
	return &reflectCache{
		nextTypes: make(map[reflect.Type]*reflectCache),
	}
}

func (rc *reflectCache) getTypeCache(objs, rest []interface{}) (*knownTypeCache, error) {
	if len(objs) == 0 {
		return nil, errors.New("objs len is 0")
	}
	if len(rest) == 0 {
		rest = objs
	}

	obj := rest[0]
	v, err := toBareStructValue(obj)
	if err != nil {
		return nil, err
	}

	if len(rest) == 1 {
		rc.RLock()
		if rc.typeCache != nil {
			rc.RUnlock()
			return rc.typeCache, nil
		}
		rc.RUnlock()

		tc, err := newTypeCache(objs)
		if err != nil {
			return nil, err
		}

		rc.Lock()
		rc.typeCache = tc
		rc.Unlock()

		return rc.typeCache, nil
	}

	v, err = toBareStructValue(rest[1])
	if err != nil {
		return nil, err
	}
	rc.RLock()
	next, ok := rc.nextTypes[v.Type()]
	rc.RUnlock()
	if !ok {
		rc.Lock()
		next = newReflectCache()
		rc.nextTypes[v.Type()] = next
		rc.Unlock()
	}

	return next.getTypeCache(objs, rest[1:])
}

func (ktc *knownTypeCache) mergeObjects(objs []interface{}) (interface{}, error) {
	s := reflect.New(ktc.st).Elem()
	for _, fa := range ktc.fas {
		obj := objs[fa.objIndex]

		v, err := toBareStructValue(obj)
		if err != nil {
			return nil, err
		}

		s.Field(fa.toFieldIndex).Set(v.Field(fa.fieldIndex))
	}

	return s.Interface(), nil
}

func newTypeCache(objs []interface{}) (*knownTypeCache, error) {
	var sfs []reflect.StructField
	var fas []*fieldAccessor
	var currentLoop int
	for objIdx, obj := range objs {
		v, err := toBareStructValue(obj)
		if err != nil {
			return nil, err
		}

		for i := 0; i < v.NumField(); i++ {
			sf := v.Type().Field(i)
			v := v.Field(i)

			if !v.CanSet() {
				continue
			}

			sfs = append(sfs, sf)
			fas = append(fas, &fieldAccessor{
				toFieldIndex: currentLoop,
				objIndex:     objIdx,
				fieldIndex:   i,
			})
			currentLoop++
		}
	}

	ktc := &knownTypeCache{
		st:  reflect.StructOf(sfs),
		fas: fas,
	}

	return ktc, nil
}

func toBareStructValue(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("must be a pointer")
	}
	if v.IsNil() {
		return reflect.Value{}, errors.New("must be a non-nil pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("must be a pointer to struct")
	}

	return v, nil
}
