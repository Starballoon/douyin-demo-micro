package dal

import (
	"context"
	"douyin-demo-micro/util"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EmailAccount  string `gorm:"column:email_account;index:unique_user,priority:8"`
	EmailDomain   string `gorm:"column:email_domain;index:unique_user,priority:9"`
	Password      []byte `gorm:"column:password;type:BINARY(64)"`
	Name          string `gorm:"column:name"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
}

func (u *User) TableName() string {
	return util.UserTable
}

func CreateUser(ctx context.Context, user *User) error {
	return DB.WithContext(ctx).Create(user).Error
}

func QueryUserByAccount(ctx context.Context, email []string) (*User, error) {
	if len(email) < 2 {
		return nil, util.ErrIllegalArguments
	}
	var user User
	if err := DB.WithContext(ctx).
		Where("email_account=? AND email_domain=?", email[0], email[1]).
		Find(&user).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return &user, nil
}

func FindUser(ctx context.Context, userID int64) (*User, error) {
	var user User
	if err := DB.WithContext(ctx).Find(&user, userID).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return &user, nil
}

func MGetUser(ctx context.Context, userIDs []int64) ([]*User, error) {
	users := make([]*User, 0)
	if len(userIDs) == 0 {
		return users, nil
	}

	if err := DB.WithContext(ctx).Find(&users, userIDs).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateDelta 这个事务的范围是否要放大
func UpdateDelta(ctx context.Context, toUserID, userID, delta int64) error {
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var users []User
		if err := tx.Find(&users, []int64{toUserID, userID}).Error; err != nil {
			return err
		}
		for _, user := range users {
			if int64(user.ID) == toUserID {
				if err := tx.Model(&user).Update("follower_count", util.MaxI64(user.FollowerCount+delta, 0)).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&user).Update("follow_count", util.MaxI64(user.FollowCount+delta, 0)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
