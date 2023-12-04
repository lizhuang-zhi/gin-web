package main

import (
	"fmt"
	"time"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).In(time.UTC).Format(TimeFormat)
}

func FormatTimestampToDate(timestamp int64) string {
	return time.Unix(timestamp, 0).In(time.UTC).Format(DateFormat)
}

func FormatTimestampLoc(timestamp int64, loc *time.Location) string {
	return time.Unix(timestamp, 0).In(loc).Format(TimeFormat)
}

func main() {
	timestamp := time.Now().Unix()
	println(timestamp)                        // 1701652843
	println(FormatTimestamp(timestamp))       // 2023-12-04 01:20:43
	println(FormatTimestampToDate(timestamp)) // 2023-12-04

	// 采用上海时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(FormatTimestampLoc(timestamp, loc)) // 2023-12-04 09:51:34

	// 解析字符串为Time类型
	parseTime, _ := time.ParseInLocation(TimeFormat, "2023-11-30 00:00:00", loc)
	fmt.Println(parseTime) // 2023-11-30 00:00:00 +0800 CST
}
