package repository

import (
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Post struct {
	Id         int64     `gorm:"column:id"`
	ParentId   int64     `gorm:"column:parent_id"`
	UserId     int64     `gorm:"column:user_id"`
	Content    string    `gorm:"column:content"`
	Diggcount  int32     `gorm:"columndigg_count"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Post) TableName() string {
	return "post"
}

type PostDao struct {
}

var postDao *PostDao
var postOnce sync.Once

func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}

func (*PostDao) QueryByPostId(id int64) (*Post, error) {
	post := &Post{}
	err := db.Where("id = ?", id).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		log.Println("QueryByPostId err")
		return nil, err
	}
	if err != nil {
		log.Println("QueryByPostId err")
		return nil, err
	} else {
		return post, nil
	}
}

func (*PostDao) QueryPostByParentId(parentid int64) ([]*Post, error) {
	var posts []*Post
	err := db.Where("parent_id = ?", parentid).Find(&posts).Error
	if err != nil {
		return nil, err
		log.Println("QueryPostByParentId err")
	}
	return posts, nil
}

func (*PostDao) CreatePost(post *Post) error {
	if err := db.Create(post).Error; err != nil {
		return err
	}
	return nil
}
