package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ArrayContains(arr interface{}, item interface{}) bool {
	newArr, ok := arr.([]interface{})
	if !ok {
		return false
	}

	for _, v := range newArr {
		if v == item {
			return true
		}
	}
	return false
}

func PrettyJson(data interface{}) string {
	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("<failed to parse json: %v>", err.Error())
	}
	return string(res)
}

func GetStructAttributesJson(s interface{}, exclude []string, excludeJsonValue []string) []string {
	// Dereference pointer if input is a pointer
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	var attributes []string

	// Loop through the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if exclude != nil {
			contain := false
			for _, v := range exclude {
				if v == field.Name {
					contain = true
					break
				}
			}
			if contain {
				continue
			}
		}

		jsonValue := field.Tag.Get("json")
		if jsonValue == "" {
			continue
		}

		if excludeJsonValue != nil {
			contain := false
			for _, v := range excludeJsonValue {
				if v == jsonValue {
					contain = true
					break
				}
			}
			if contain {
				continue
			}
		}
		attributes = append(attributes, jsonValue)
	}

	return attributes
}
