package main

import "fmt"

type User struct {
    Name   string      `json:"name"`
}

// // 定义结构体的方法
// func (u User) SetNameValue(val string) {
//     u.Name = val
// }

// 定义结构体指针的方法
func (u *User) SetNameValue(val string) {
    u.Name = val
}

func main() {
    user1 := User{
        Name:   "Leo",
    }
	user2 := User{
        Name:   "Jim",
    }
    fmt.Println(user1)  // {Leo}
    fmt.Println(user2)  // {Jim}

    user1.SetNameValue("Jeo")
    fmt.Println(user1)  // {Jeo}
    fmt.Println(user2)  // {Jim}

    user2.SetNameValue("Tim")
    fmt.Println(user1)  // {Jeo}
    fmt.Println(user2)  // {Tim}
}