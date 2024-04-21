package main

import (
	"fmt"
	"reflect"
)

type People struct {
	Name      string      `json:"name" isloc:"true"`
	Age       int         `json:"age"`
	LikeSport *Basketball `json:"like_sport"`
	LikeMusic *Music      `json:"like_music"`
}

type Basketball struct {
	Title string `json:"title" isloc:"true"`
	Count int    `json:"count"`
}

type Music struct {
	MusicName string `json:"music_name" isloc:"true"`
	MusicType string `json:"music_type"`
}

var localMap []string

func main() {
	// 构建测试数据
	likeSport := &Basketball{
		Title: "100",
		Count: 1,
	}
	likeMusic := &Music{
		MusicName: "200",
		MusicType: "pop",
	}
	people := People{
		Name:      "300",
		Age:       3,
		LikeSport: likeSport,
		LikeMusic: likeMusic,
	}

	// 通过反射检查 People 结构体中的本地化字段
	TraverseObject(&people)

	fmt.Println(localMap)
}

func TraverseObject(obj interface{}) {
	peopleVal := reflect.ValueOf(obj)

	if peopleVal.Kind() == reflect.Ptr {
		peopleVal = peopleVal.Elem()
	}

	if peopleVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < peopleVal.NumField(); i++ {
		fieldVal := peopleVal.Field(i)
		fieldType := peopleVal.Type().Field(i)

		if fieldVal.Kind() == reflect.Struct {
			TraverseObject(fieldVal.Interface())
			continue
		}

		if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
			fieldVal = fieldVal.Elem()
			TraverseObject(fieldVal.Interface())
			continue
		}

		if fieldType.Tag.Get("isloc") == "true" {
			localMap = append(localMap, fieldVal.String())
		}
	}
}
