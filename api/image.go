package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
	"log"
)

func PullImage(ctx context.Context, cli *client.Client, name string) ([]jsonmessage.JSONMessage, error) {
	stream, err := cli.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	decoder := json.NewDecoder(stream)
	var allMessages []jsonmessage.JSONMessage
	var pullSuccess = true

	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				// 到达流的末尾，正常结束
				break
			}
			// 解析错误应该返回，而不是致命退出
			return nil, fmt.Errorf("解析消息失败: %v", err)
		}
		// 保存消息
		allMessages = append(allMessages, msg)
		// 检查错误
		if msg.Error != nil {
			pullSuccess = false
			// 记录错误但继续处理其他消息，不要直接退出
			log.Printf("拉取错误: %v", msg.Error)
		}
	}
	// 拉取完成后，检查结果
	if !pullSuccess {
		return nil, fmt.Errorf("镜像拉取过程中发生错误，请检查日志")
	}
	return allMessages, nil
}

func Rmi(ctx context.Context, cli *client.Client, imageName string) ([]image.DeleteResponse, error) {
	return cli.ImageRemove(ctx, imageName, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
}
