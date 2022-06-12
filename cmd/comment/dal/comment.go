package dal

import (
	"context"
	"douyin-demo-micro/util"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	VideoID    int64  `gorm:"column:video_id;index"`
	ReviewerID int64  `gorm:"column:reviewer_id;"`
	Content    string `gorm:"column:content"`
}

func (u *Comment) TableName() string {
	return util.CommentTable
}

func CreateComment(ctx context.Context, comment *Comment) error {
	if err := DB.WithContext(ctx).Create(comment).Error; err != nil {
		return util.ErrInternalError
	}
	return nil
}

func DeleteComment(ctx context.Context, comment *Comment) error {
	if err := DB.WithContext(ctx).Delete(comment).Error; err != nil {
		return util.ErrInternalError
	}
	return nil
}

func FindComment(ctx context.Context, commentID int64) (*Comment, error) {
	var comment Comment
	if err := DB.WithContext(ctx).Find(&comment, commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func CommentListByVideoID(ctx context.Context, videoID int64) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("video_id=?", videoID).Order("created_at DESC").
		Limit(100).Find(&comments).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return comments, nil
}
