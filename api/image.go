package api

import (
	"context"
	"docker-mcp/cmd/logs"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"io"
)

func PullImage(ctx context.Context, cli *client.Client, name string) ([]jsonmessage.JSONMessage, error) {
	logs.InfoWithFields("Start pulling image", map[string]interface{}{"image": name})
	stream, err := cli.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		logs.ErrorWithFields("ImagePull failed", map[string]interface{}{"image": name, "error": err})
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
				logs.InfoWithFields("Image pull stream EOF", map[string]interface{}{"image": name})
				break
			}
			logs.ErrorWithFields("Decode image pull message failed", map[string]interface{}{"image": name, "error": err})
			return nil, fmt.Errorf("解析消息失败: %v", err)
		}
		// 保存消息
		allMessages = append(allMessages, msg)
		// 检查错误
		if msg.Error != nil {
			pullSuccess = false
			logs.ErrorWithFields("Image pull error message", map[string]interface{}{"image": name, "error": msg.Error})
		}
		logs.DebugWithFields("Image pull message", map[string]interface{}{"image": name, "msg": msg})
	}
	// 拉取完成后，检查结果
	if !pullSuccess {
		logs.ErrorWithFields("Image pull finished with error", map[string]interface{}{"image": name})
		return nil, fmt.Errorf("镜像拉取过程中发生错误，请检查日志")
	}
	logs.InfoWithFields("Image pull finished successfully", map[string]interface{}{"image": name})
	return allMessages, nil
}

func Rmi(ctx context.Context, cli *client.Client, imageName string) ([]image.DeleteResponse, error) {
	logs.InfoWithFields("Start removing image", map[string]interface{}{"image": imageName})
	resp, err := cli.ImageRemove(ctx, imageName, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	if err != nil {
		logs.ErrorWithFields("Remove image failed", map[string]interface{}{"image": imageName, "error": err})
		return nil, err
	}
	logs.InfoWithFields("Remove image success", map[string]interface{}{"image": imageName, "result": resp})
	return resp, nil
}
