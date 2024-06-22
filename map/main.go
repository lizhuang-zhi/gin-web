package main

import (
	"fmt"
	"time"
)

var GlobalMap map[string]string
var InitSlices = []map[string]string{
	{"a": "leo"},
	{"b": "tim"},
	{"c": "jeo"},
}

func init() {
	// GlobalMap = make(map[string]string)  // 这里不能初始化，否则会导致后续更新数据时，map的数据不会被清空ß
	updateMap(InitSlices)
}

func updateMap(slices []map[string]string) {
	GlobalMap = make(map[string]string, len(slices)) // 重置map

	for _, val := range slices {
		for k, v := range val {
			GlobalMap[k] = v
		}
	}
}

func main() {
	fmt.Println("初始数据:", GlobalMap)

	time.Sleep(1 * time.Second)

	// 模拟更新数据
	newSlices := []map[string]string{
		{"a": "leo111"},
		{"e": "jim"},
	}

	updateMap(newSlices)

	fmt.Println("更新后:", GlobalMap)
}
