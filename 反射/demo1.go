package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4

	// 使用 reflect.TypeOf() 获取变量的类型
	fmt.Println("type:", reflect.TypeOf(x))

	// 使用 reflect.ValueOf() 获取变量的值
	v := reflect.ValueOf(x)
	fmt.Println("value:", v)
	/*
		在Go语言的反射中，Type 和 Kind 两者虽然联系紧密，但意义不同：
		Type 表示接口中存储的具体类型，它是 reflect.Type 类型的对象，提供了描述类型的详细信息，例如类型的名称、是否为指针类型、是否实现了某个接口等。Type 可以让你知道一个变量是 int、float64、MyStruct 这类具体的类型。
		Kind 也是一个表示类型的值，但它描述的是基础分类，或者说类型的底层种类。所有基础的Go类型（如整数、浮点数、复数、数组、结构体、指针等）都是 reflect.Kind 类型的值。Kind 用于区分一个类型属于哪个基本种类，因此比起 Type，Kind 更加抽象，它是通过 reflect.Kind 常量来表示的。
		这意味着不同的 Type 可能对应相同的 Kind。例如，不论是 int32 还是 int64，它们的 Kind 都是 reflect.Int，但它们的 Type 是不同的。

		type MyInt int
		var x MyInt = 7

		v := reflect.ValueOf(x)
		fmt.Println("type:", v.Type()) // 输出自定义类型 MyInt
		fmt.Println("kind:", v.Kind()) // 输出基础种类 int
	*/
	fmt.Println("type:", v.Type())
	fmt.Println("kind:", v.Kind())

	// 使用反射修改变量的值
	// 需要确保使用 ValueOf 方法时传入变量的指针
	p := reflect.ValueOf(&x) // 注意这里传入的是 x 的地址
	v = p.Elem()
	fmt.Println("settability of v:", v.CanSet())

	v.SetFloat(7.1)
	fmt.Println("x:", x)
	fmt.Println("v:", v)
}
