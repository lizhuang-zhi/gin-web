package context

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	context.Context
	context.Context 是 Golang 标准库中的一个包，用于控制 goroutine 生命周期和传递请求范围的数据。它主要用于以下场景：

	取消操作：传递取消信号，通知 goroutine 停止操作。
	超时控制：设置操作的截止时间。
	传递数据：在请求范围内传递少量关键数据（如认证信息、请求 ID 等）。
	context 包主要有四种类型的上下文：

	context.Background()：返回一个空的上下文，通常用于主函数、初始化或测试。
	context.TODO()：返回一个空的上下文，表示当前还不知道使用哪种上下文。
	context.WithCancel(parent)：返回一个带取消功能的上下文及取消函数。
	context.WithDeadline(parent, deadline) 和 context.WithTimeout(parent, timeout)：返回一个带截止时间或超时功能的上下文。


	gin.Context
	gin.Context 是 Gin Web 框架中的上下文，用于在处理 HTTP 请求时传递各种信息和操作方法。gin.Context 提供了以下功能：

	请求处理：获取请求的路径、查询参数、表单数据等。
	响应处理：设置响应状态码、响应头、响应数据等。
	中间件：在请求处理链中传递数据和控制流。
	错误处理：记录和处理请求过程中的错误。

	自我总结: gin.Context应该是单个请求中的关联上下文,也就是每个单次的http请求中的c *gin.Context,
	就是贯穿从获取请求参数, 到去数据库(model层)中查询,再到返回响应数据的这一个过程中的上下文记录
*/

/*
	源代码:
*/
// 添加cluster信息
func AddCluster(ctx *gin.Context) (interface{}, error) {
	params := &AddClusterParam{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	id, err := core.ClusterUUID.Get()
	if err != nil {
		return nil, err
	}
	cluster := &model.Cluster{
		ClusterID:  int(id),
		Name:       params.Name,
		APIURL:     params.APIURL,
		UpdateTime: utils.NowFormat(),
	}
	info, err := sdk.GetClusterInfo(params.APIURL)
	if err != nil {
		return nil, err
	}
	cluster.ClusterInfo = info
	return nil, model.UpsertCluster(cluster)
}

/*
上面的源代码,此时在传递到下一层model时,并没有设置context,进行修改
这样改动后, 便可以对本次的http请求进行超时控制!
*/
func AddCluster(ctx *gin.Context) (interface{}, error) {
	// ....

	// 设置超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// ....

	// 传递带超时的 context
	return nil, model.UpsertCluster(timeoutCtx, cluster)
}
