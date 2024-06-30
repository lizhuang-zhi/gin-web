package commandpb

// TODO: 此文件应该是通过模板自动生成(这里只是为了方便)

type Command int32

const (
	Command_None Command = 0
	/*
		网关服务下的grpc方法(1000-10000)
	*/
	Command_GateLogin Command = 1000 // 登录
	Command_GatePing  Command = 1001 // Ping
	/*
		活动服务下的grpc方法(10000-20000)
	*/
	Command_ActivityGetNotice        Command = 10000 // 获取活动公告
	Command_ActivityGetNotices       Command = 10001 // 获取所有活动公告
	Command_ActivityUpdateNoticeById Command = 10001 // 通过id更新活动公告
	/*
		邮件服务下的grpc方法(20000-30000)
	*/
	Command_MailGetAllMail Command = 20000 // 获取所有邮件
)
