package main

import (
	"context"
	"fmt"

	sdk "github.com/larksuite/project-oapi-sdk-golang"
	sdkcore "github.com/larksuite/project-oapi-sdk-golang/core"
	"github.com/larksuite/project-oapi-sdk-golang/service/project"
)

func main() {
	// 创建 client
	client := sdk.NewClient("MII_6692498323C08001", "940917DB647B3CD7D13E8540520CD65A")

	// 发起请求
	resp, err := client.Project.ListProjectWorkItemType(context.Background(), project.NewListProjectWorkItemTypeReqBuilder().
		ProjectKey("64fc1eecb1232854aca226a2").
		Build(),
		sdkcore.WithUserKey("7234833414560546819"),
	)

	//处理错误
	if err != nil {
		// 处理err
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code(), resp.ErrMsg, resp.RequestId())
		return
	}

	// 业务数据处理
	fmt.Println(sdkcore.Prettify(resp.Data))
}
