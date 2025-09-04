package utils

import (
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// опускает пустые значения
func StructToBson(st any) bson.D {
	val := reflect.ValueOf(st)
	m := bson.D{}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return m
	}
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		if fieldVal.IsValid() {
			fieldName := t.Field(i)
			m = append(m, bson.E{Key: fieldName.Name, Value: fieldVal.Interface()})
		}
	}

	return m
}
