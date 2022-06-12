package service

import (
	"douyin-demo-micro/cmd/user/dal"
	"douyin-demo-micro/kitex_gen/user"
)

func User(u *dal.User) *user.User {
	if u == nil {
		return nil
	}
	return &user.User{
		Id:            int64(u.ID),
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
	}
}

func Users(us []*dal.User) []*user.User {
	users := make([]*user.User, 0)
	for _, u := range us {
		if u2 := User(u); u2 != nil {
			users = append(users, u2)
		}
	}
	return users
}
