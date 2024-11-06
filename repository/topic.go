package repository

import (
	"log"
	"sync"
	"time"
)

type Topic struct {
	Id         int64     `gorm:"column:id"`
	UserId     int64     `gorm:"column:user_id"`
	Title      string    `gorm:"column:title"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Topic) TableName() string {
	return "topic"
}

type TopicDao struct {
}

var topicDao *TopicDao
var topicOnce sync.Once

func NewTopicDaoInstance() *TopicDao {
	topicOnce.Do(
		func() {
			topicDao = &TopicDao{}
		})
	return topicDao
}

func (*TopicDao) QueryByTopicId(id int64) (*Topic, error) {
	topic := &Topic{}
	if err := db.Where("id = ?", id).First(topic).Error; err != nil {
		log.Println("QueryByTopicId err")
		return nil, err
	}
	return topic, nil
}
