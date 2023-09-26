package service

import (
	"context"
	"fmt"
)

type IPostService interface {
	GetPostById(ctx context.Context, id string) string
}

//var _ IPostService = (*PostService)(nil)

type PostService struct {
}

func (s *PostService) GetPostById(ctx context.Context, id string) string {
	return fmt.Sprint("欢迎关注公众号：Go技术干货，作者：陈明勇")
}

// 使用这个函数，则wire处不需要进行改造
func NewPostService() IPostService {
	return &PostService{}
}

// 使用这个函数，需要通过wire.Bind 建立IPostService 和 *PostService的关系
func NewPostServiceV3() *PostService {
	return &PostService{}
}
