package main

import "fmt"

type Option func(*Person) error

// *********************** Person ***********************
// 携带所有服务需要用到的资源
type Person struct {
	// 基本信息
	Name string
	Age  int
	City string

	MusicBox bool // 是否有音乐盒(第三方工具)
	Mirror   bool // 是否有镜子
	Handrail bool // 是否有舞蹈室栏杆

	Computer bool // 是否有电脑

	Pot    bool // 是否有锅
	Shovel bool // 是否有铲子
}

func NewPerson() *Person {
	return &Person{}
}

func (p *Person) WithBaseInfo(name string, age int, city string) Option {
	return func(p *Person) error {
		p.Name = name
		p.Age = age
		p.City = city
		return nil
	}
}

func (p *Person) WithMusicBox(musicBox bool) Option {
	return func(p *Person) error {
		p.MusicBox = musicBox
		return nil
	}
}

func (p *Person) WithMirror(mirror bool) Option {
	return func(p *Person) error {
		p.Mirror = mirror
		return nil
	}
}

func (p *Person) WithComputer(computer bool) Option {
	return func(p *Person) error {
		p.Computer = computer
		return nil
	}
}

// *********************** 跳舞服务 ***********************
type DanceService struct {
	Person *Person
}

func NewDanceService() *DanceService {
	return &DanceService{
		Person: NewPerson(),
	}
}

func (d *DanceService) Start(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(d.Person); err != nil {
			return err
		}
	}

	return nil
}

// *********************** 开发服务 ***********************
type DevelopService struct {
	Person *Person
}

func NewDevelopService() *DevelopService {
	return &DevelopService{
		Person: NewPerson(),
	}
}

func (d *DevelopService) Start(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(d.Person); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// 跳舞服务
	danceService := NewDanceService()
	err := danceService.Start(
		danceService.Person.WithBaseInfo("张三", 18, "北京"), // 基本信息(姓名、年龄、城市
		danceService.Person.WithMusicBox(true),           // 启动音乐盒插件
		danceService.Person.WithMirror(true),             // 启动镜子插件
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("跳舞服务启动成功！")

	fmt.Println(danceService.Person.Name)     // 姓名: 张三
	fmt.Println(danceService.Person.Age)      // 年龄: 18
	fmt.Println(danceService.Person.MusicBox) // MusicBox启动成功: true
	fmt.Println(danceService.Person.Computer) // Computer插件未启动: false

	// 开发服务
	developService := NewDevelopService()
	err = developService.Start(
		developService.Person.WithBaseInfo("李四", 22, "上海"), // 基本信息(姓名、年龄、城市
		developService.Person.WithComputer(true),           // 启动电脑插件
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("开发服务启动成功！")

	fmt.Println(developService.Person.Name)     // 姓名: 李四
	fmt.Println(developService.Person.Age)      // 年龄: 22
	fmt.Println(developService.Person.Computer) // Computer插件未启动: true
	fmt.Println(developService.Person.MusicBox) // MusicBox启动成功: false
}
