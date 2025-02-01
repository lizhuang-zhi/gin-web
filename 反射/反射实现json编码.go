// 实现一个简单的 JSON 编码器，使用反射将结构体转换为 JSON 字符串（忽略嵌套和复杂类型）。
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

func main() {
	user := User{ // 这里不能是指针
		Name: "leo",
		Age:  24,
		Job:  "coder",
	}

	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)

	fieldList := make([]string, 0)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("json")
		key := fieldType.Name
		if tag != "" {
			key = tag
		}

		value := ""
		switch field.Kind() {
		case reflect.String:
			value = fmt.Sprintf(`%s`, field.String())
		case reflect.Int:
			value = fmt.Sprintf(`%d`, field.Int())
		}

		fieldList = append(fieldList, fmt.Sprintf(`"%s":%s`, key, value))
	}

	result := "{" + strings.Join(fieldList, ",") + "}"
	fmt.Println(result)
}
