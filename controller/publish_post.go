package controller

import (
	"strconv"
	"youthcamp/lesson02/project/service"
)

func PublishPost(uidStr string, topicIdStr string, content string) *PageData {
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	topicId, _ := strconv.ParseInt(topicIdStr, 10, 64)
	postId, err := service.PublishPost(topicId, uid, content)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PageData{
		Code: 0,
		Msg:  "success",
		Data: map[string]int64{
			"post_id": postId,
		},
	}

}
