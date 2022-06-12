package service

import (
	"bytes"
	"context"
	"douyin-demo-micro/cmd/user/dal"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
	"strings"
)

func CreateUser(ctx context.Context, req *user.CreateUserRequest) error {
	email := strings.Split(req.Req.Username, "@")
	user, err := dal.QueryUserByAccount(ctx, email)
	if err != nil {
		return err
	}
	if user.ID > 0 {
		return util.ErrUserAlreadyExist
	}

	err = dal.CreateUser(ctx, &dal.User{
		EmailAccount: email[0],
		EmailDomain:  email[1],
		Password:     util.Encrypt(req.Req.Password),
		// TODO 新用户名生成策略
		Name: "NewUser",
	})
	if err != nil {
		return util.ErrInternalError
	}
	return nil
}

func CheckUser(ctx context.Context, req *user.CheckUserRequest) (*user.User, error) {
	email := strings.Split(req.Req.Username, "@")
	user, err := dal.QueryUserByAccount(ctx, email)
	if err != nil {
		return nil, err
	}
	if user.ID > 0 {
		password := util.Encrypt(req.Req.Password)
		if bytes.Equal(user.Password, password) {
			return User(user), nil
		}
	}
	return nil, util.ErrUserNotFound
}

func FindUser(ctx context.Context, req *user.FindUserRequest) (*user.User, error) {
	if req.Req.UserId == 0 {
		return nil, util.ErrIllegalArguments
	}

	user, err := dal.FindUser(ctx, req.Req.UserId)
	if err != nil {
		return nil, err
	}
	return User(user), nil
}

func MGetUser(ctx context.Context, req *user.MGetUserRequest) ([]*user.User, error) {
	if len(req.UserIds) == 0 {
		return []*user.User{}, nil
	}

	users, err := dal.MGetUser(ctx, req.UserIds)
	if err != nil {
		return nil, util.ErrInternalError
	}

	if req.UserId == nil && *req.UserId > 0 {
		result, err := merge(ctx, users, *req.UserId)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return Users(users), nil
}

func merge(ctx context.Context, users []*dal.User, userID int64) ([]*user.User, error) {
	leaderIDs, err := dal.QueryLeaderIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	leaderIDMap := make(map[int64]struct{})
	for _, id := range leaderIDs {
		leaderIDMap[id] = struct{}{}
	}

	result := Users(users)
	for _, v := range result {
		_, v.IsFollow = leaderIDMap[v.Id]
	}
	return result, nil
}
