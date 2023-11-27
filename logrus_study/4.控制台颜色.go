package main

import (
	"fmt"
)

const (
	grayBlock   = 0
	redBlock    = 1
	greenBlock  = 2
	yellowBlock = 3
)

func PrintColor(colorCode int, text string, isBackground bool) {
	if !isBackground {
		fmt.Printf("\033[3%dm %s \033[0m\n", colorCode, text)
		return 
	}
	fmt.Printf("\033[4%dm %s \033[0m\n", colorCode, text)
}

func main() {
	// 不同输出颜色(字体颜色)
	fmt.Println("\033[30m 灰色 \033[0m")
	fmt.Println("\033[31m 红色 \033[0m")
	fmt.Println("\033[32m 绿色 \033[0m")
	fmt.Println("\033[33m 黄色 \033[0m")

	// 不同输出颜色(背景色)
	fmt.Println("\033[40m 灰色 \033[0m")
	fmt.Println("\033[41m 红色 \033[0m")
	fmt.Println("\033[42m 绿色 \033[0m")
	fmt.Println("\033[43m 黄色 \033[0m")

	PrintColor(greenBlock, "测试输出", false)
	PrintColor(yellowBlock, "测试输出", false)
	PrintColor(grayBlock, "测试输出", true)
	PrintColor(redBlock, "测试输出", true)
}
