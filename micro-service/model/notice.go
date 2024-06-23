package model

// 公告数据
type Notice struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	SubType int    `json:"sub_type"`
	Content string `json:"content"`
}

// 内存中的公告数据
var GlobalNotice = make(map[int]*Notice)
