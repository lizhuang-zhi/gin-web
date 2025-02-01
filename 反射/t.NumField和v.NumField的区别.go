// 自我总结：针对结构体使用时，无太大差异

/*
	v.NumField() 和 t.NumField() 的区别
	v.NumField()：
	属于 reflect.Value 的方法，直接基于 值的实际类型 返回字段数量。

	如果 v 是无效的 reflect.Value（例如未初始化或空接口），调用 v.NumField() 会 panic。

	如果 v 的类型不是 struct，调用 v.NumField() 也会 panic。

	t.NumField()：
	属于 reflect.Type 的方法，基于 类型的元信息 返回字段数量。

	如果 t 的类型不是 struct，调用 t.NumField() 会 panic。
*/

/*
例子一：v.NumField() 的底层实现等同于 v.Type().NumField()，即先通过 v.Type() 获取类型信息，再调用 NumField()。
因此，当 v 是结构体值时，两者的结果相同。

type User struct {
	Name string
	Age  int
}

user := User{"Alice", 30}

t := reflect.TypeOf(user)  // 类型是 User
v := reflect.ValueOf(user) // 值是 User 的实例

fmt.Println(t.NumField()) // 输出: 2
fmt.Println(v.NumField()) // 输出: 2
*/

/*
例子二：如果 v 是一个指针，v.NumField() 会直接 panic（因为指针类型没有字段），但 t.NumField() 的行为取决于 t 的类型

userPtr := &User{"Bob", 25}

t := reflect.TypeOf(userPtr)  // 类型是 *User
v := reflect.ValueOf(userPtr) // 值是指针

fmt.Println(t.NumField())     // panic: 类型是 *User（指针），没有字段！
fmt.Println(v.NumField())     // panic: 值的类型是 *User（指针），没有字段！

// 正确做法：解引用指针
v = v.Elem()                  // 获取指针指向的值（User 实例）
t = t.Elem()                  // 获取指针指向的类型（User 类型）
fmt.Println(t.NumField())     // 输出: 2
fmt.Println(v.NumField())     // 输出: 2
*/