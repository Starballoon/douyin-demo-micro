package service

import (
	"douyin-demo-micro/cmd/comment/dal"
	"douyin-demo-micro/kitex_gen/comment"
	"douyin-demo-micro/kitex_gen/user"
)

func Comment(c *dal.Comment, u *user.User) *comment.Comment {
	return &comment.Comment{
		Id:         int64(c.ID),
		User:       u,
		Content:    c.Content,
		CreateDate: c.CreatedAt.Format("0102"),
	}
}

func Comments(cs []*dal.Comment, us []*user.User) []*comment.Comment {
	reviewerMap := make(map[int64]*user.User)
	for _, reviewer := range us {
		reviewerMap[reviewer.Id] = reviewer
	}

	var results []*comment.Comment
	for _, c := range cs {
		reviewer, _ := reviewerMap[c.ReviewerID]
		results = append(results, Comment(c, reviewer))
	}
	return results
}
