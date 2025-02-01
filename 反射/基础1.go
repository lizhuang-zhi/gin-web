package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 构建测试数据
	person := Person{
		Name: "300",
		Age:  3,
	}

	v := reflect.ValueOf(person)
	t := reflect.TypeOf(person)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf("field name: %s, field type: %s, field value: %v\n", fieldType.Name, fieldType.Type, field.Interface())
	}
}
