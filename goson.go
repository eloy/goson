package goson

import (
	"strings"
	"reflect"
	"encoding/json"
)

type Goson interface {
	ToJson() ([]byte, error)
	Alias(string, string) Goson
	Method(string) Goson
	Hash(string, ...string) Goson
	HashAlias(string, string, ...string) Goson
}

type goson struct {
	methods []string
	alias map[string]string
	hashes []*GosonNestedHash
}

func (this *goson) initialize() {
	this.alias = make(map[string]string)
	this.hashes = make([]*GosonNestedHash, 0)
}

type Values map[string]interface{}

type GosonHash struct {
	goson
	model interface{}
}

type GosonNestedHash struct {
	goson
	method string
	alias string
}

type GosonArray struct {
	goson
	size int
	array interface{}
}


func Hash(model interface{}, methods...string) *GosonHash {
	h := new(GosonHash)
	h.goson.initialize()
	h.model = model
	h.methods = methods
	return h
}


func Array(data interface{}, methods...string) *GosonArray {
	a := new(GosonArray)
	a.goson.initialize()
	a.array = data
	a.methods = methods
	return a
}

func (this *goson) Alias(alias string, method string) Goson {
	this.alias[alias] = method
	return this
}

func (this *goson) Method(method string) Goson {
	this.methods = append(this.methods, method)
	return this
}

func (this *goson) Hash(key string, methods...string) Goson {
	nestedHash := new(GosonNestedHash)
	nestedHash.goson.initialize()
	nestedHash.method = key
	nestedHash.alias = key
	nestedHash.methods = methods

	this.hashes = append(this.hashes, nestedHash)
	return nestedHash
}

func (this *goson) HashAlias(key string, alias string, methods...string) Goson {
	nestedHash := new(GosonNestedHash)
	nestedHash.goson.initialize()
	nestedHash.method = key
	nestedHash.alias = alias
	nestedHash.methods = methods

	this.hashes = append(this.hashes, nestedHash)
	return nestedHash
}


func (this *goson) ToJson() ([]byte,error) {
	panic("should never happen")
}

func (this *GosonHash) ToJson() ([]byte,error) {
	hash := this.toMap(this.model)
	return json.Marshal(hash)
}


func (this *GosonArray) ToJson() ([]byte,error) {
	slice := reflect.ValueOf(this.array)
	length := slice.Len()
	array := make([]Values, length)

	for i := 0; i < length; i++ {
		array[i] = this.toMap(slice.Index(i).Interface())
	}

	return json.Marshal(array)
}

func (this *goson) toMap(model interface{}) Values {
	hash := make(Values)
	value := reflect.ValueOf(model)

	value = reflect.Indirect(value)

	for _, name := range this.methods {
		hash[prettyName(name)] = getValue(name, value)
	}

	for alias, name := range this.alias {
		hash[prettyName(alias)] = getValue(name, value)
	}

	for _, h := range this.hashes {
		model := getValue(h.method, value)
		if model != nil {
			hash[prettyName(h.alias)] = h.toMap(model)
		}
	}

	return hash
}


func getValue(name string, value reflect.Value) interface{} {
	if strings.HasSuffix(name, "()") {
		methodName := name[:len(name) - 2]
		return value.MethodByName(methodName).Call(nil)[0].Interface()
	} else {
		return value.FieldByName(name).Interface()
	}
}

func prettyName(name string) string {
	name = strings.ToLower(name)
	if strings.HasSuffix(name, "()") {
		name = name[:len(name) - 2]
	}
	return name
}
