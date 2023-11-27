package main

import "fmt"

type User struct {
    UserID int         `json:"user_id"`
    Name   string      `json:"name"`
    Age    int         `json:"age"`
    Wallet interface{} `json:"wallet"`
}

// 定义结构体的方法
func (u User) SetNameValue(val string) {
    u.Name = val
}

// 定义结构体指针的方法
func (u *User) SetNamePointer(val string) {
    u.Name = val
}

func main() {
    user := User{
        UserID: 1,
        Name:   "Leo",
        Age:    22,
        Wallet: 10000,
    }
    fmt.Println(user)  // {1 Leo 22 10000}

    user.SetNameValue("Jim")
    // 由于是对值的拷贝, 所以不会将改变带出SetNameValue函数
    fmt.Println(user)  // {1 Leo 22 10000} 

    user.SetNamePointer("Tim")
    // 由于是对结构体地址的访问, 所以会将改变带出SetNamePointer函数
    fmt.Println(user)  // {1 Tim 22 10000}
}