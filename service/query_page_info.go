package service

import (
	"errors"
	"fmt"
	"sync"
	"youthcamp/lesson02/project/repository"
)

type TopicInfo struct {
	Topic *repository.Topic
	User  *repository.User
}

type PostInfo struct {
	Post *repository.Post
	User *repository.User
}

type PageInfo struct {
	TopicInfo *TopicInfo
	PostInfo  []*PostInfo
}

// QueryPageInfo 查询帖子详情
func QueryPageInfo(topicId int64) (*PageInfo, error) {
	return NewQueryPageInfoFlow(topicId).Do()
}

func NewQueryPageInfoFlow(topId int64) *QueryPageInfoFlow {
	return &QueryPageInfoFlow{
		topicId: topId,
	}
}

type QueryPageInfoFlow struct {
	topicId  int64
	pageInfo *PageInfo
	topic    *repository.Topic
	posts    []*repository.Post
	userMap  map[int64]*repository.User
}

func (f *QueryPageInfoFlow) Do() (*PageInfo, error) {
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	if err := f.packPageInfo(); err != nil {
		return nil, err
	}
	return f.pageInfo, nil
}

func (f *QueryPageInfoFlow) checkParam() error {
	if f.topicId <= 0 {
		return errors.New("topicId must be larger than 0")
	}
	return nil
}

// 通过topicID初始化topic，post列表和userMap
func (f *QueryPageInfoFlow) prepareInfo() error {
	var wg sync.WaitGroup
	wg.Add(2)
	var topicErr, postErr error
	go func() {
		defer wg.Done()
		topic, err := repository.NewTopicDaoInstance().QueryByTopicId(f.topicId)
		if err != nil {
			topicErr = err
			return
		}
		f.topic = topic
	}()

	go func() {
		defer wg.Done()
		posts, err := repository.NewPostDaoInstance().QueryPostByParentId(f.topicId)
		if err != nil {
			postErr = err
			return
		}
		f.posts = posts
	}()

	wg.Wait()

	if topicErr != nil {
		return topicErr
	}
	if postErr != nil {
		return postErr
	}

	uids := []int64{f.topicId}
	for _, post := range f.posts {
		uids = append(uids, post.Id)
	}
	userMap, err := repository.NewUserDaoInstance().MQueryUserById(uids)
	if err != nil {
		return err
	}
	f.userMap = userMap
	return nil
}

// /初始化pageinfo
func (f *QueryPageInfoFlow) packPageInfo() error {
	//topic info
	userMap := f.userMap
	topicUser, ok := userMap[f.topic.UserId]
	if !ok {
		return errors.New("has no topic user info")
	}
	//post list
	postList := make([]*PostInfo, 0)
	for _, post := range f.posts {
		postUser, ok := userMap[post.UserId]
		if !ok {
			return errors.New("has no post user info for " + fmt.Sprint(post.UserId))
		}
		postList = append(postList, &PostInfo{
			Post: post,
			User: postUser,
		})
	}
	f.pageInfo = &PageInfo{
		TopicInfo: &TopicInfo{
			Topic: f.topic,
			User:  topicUser,
		},
		PostInfo: postList,
	}
	return nil
}
