package dal

import (
	"context"
	"douyin-demo-micro/util"
	"gorm.io/gorm"
)

type Following struct {
	gorm.Model
	LeaderID   int64 `gorm:"column:leader_id;index:positive;index:reverse"`
	FollowerID int64 `gorm:"column:follower_id;index:positive;index:reverse,priority:9"`
}

func (u *Following) TableName() string {
	return util.FollowTable
}

func CreateFollow(ctx context.Context, follow *Following) error {
	return DB.WithContext(ctx).Create(follow).Error
}

func RecoverFollowUnscoped(ctx context.Context, follow *Following) error {
	return DB.WithContext(ctx).Unscoped().
		Model(&follow).Update("deleted_at", nil).Error
}

func QueryFollowUnscoped(ctx context.Context, toUserID, userID int64) (*Following, error) {
	var follow Following
	// Unscoped会把软删除的行也查询出来
	if err := DB.WithContext(ctx).Unscoped().
		Where("leader_id=? AND follower_id=?", toUserID, userID).Find(&follow).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return &follow, nil
}

func DeleteFollow(ctx context.Context, follow *Following) error {
	return DB.WithContext(ctx).Delete(&follow).Error
}

func QueryLeaderIDs(ctx context.Context, userID int64) ([]int64, error) {
	leaderIDs := make([]int64, 0)
	if err := DB.WithContext(ctx).Table(util.FollowTable).
		Select("lead_id").Where("follower_id=?", userID).Find(&leaderIDs).Error; err != nil {
		return nil, err
	}
	return leaderIDs, nil
}

func QueryLeaderIDByUserID(ctx context.Context, leaderIDs []int64, userID int64) ([]int64, error) {
	result := make([]int64, 0)
	if err := DB.WithContext(ctx).Table(util.FollowTable).
		Select("lead_id").Where("leader_id IN (?) AND follower_id=?", leaderIDs, userID).Find(&result).Error; err != nil {
		return nil, err
	}
	return leaderIDs, nil
}

func FollowList(ctx context.Context, userID int64) ([]*User, error) {
	leadQuery := DB.WithContext(ctx).Select("leader_id").Where("follower_id=?", userID).Table(util.FollowTable)
	var leaders []*User
	err := DB.WithContext(ctx).Where("id IN (?)", leadQuery).Limit(100).Find(&leaders).Error
	if err != nil {
		return nil, util.ErrInternalError
	}
	return leaders, nil
}

func FollowerList(ctx context.Context, userID int64) ([]*User, error) {
	leadQuery := DB.WithContext(ctx).Select("follower_id").Where("leader_id=?", userID).Table(util.FollowTable)
	var leaders []*User
	err := DB.WithContext(ctx).Where("id IN (?)", leadQuery).Limit(100).Find(&leaders).Error
	if err != nil {
		return nil, util.ErrInternalError
	}
	return leaders, nil
}
