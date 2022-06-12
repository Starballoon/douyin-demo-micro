package service

import (
	"context"
	"douyin-demo-micro/cmd/user/dal"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
)

// Follow 参数默认是合法的，即用户存在且操作合法，这个操作要4次左右的数据库查询，前端没有做重复请求拦截，这个有必要在外层优化
func Follow(ctx context.Context, req *user.FollowRequest) error {
	toUserID := req.LeaderId
	userID := req.FollowerId
	if toUserID == userID || toUserID == 0 || userID == 0 {
		return util.ErrIllegalArguments
	}

	follow, err := dal.QueryFollowUnscoped(ctx, toUserID, userID)
	if err != nil {
		return err
	}
	if follow.ID > 0 && follow.Model.DeletedAt.Valid {
		if err = dal.RecoverFollowUnscoped(ctx, follow); err != nil {
			return util.ErrInternalError
		}
		if err = dal.UpdateDelta(ctx, toUserID, userID, 1); err != nil {
			return util.ErrInternalError
		}
	} else if follow.ID == 0 {
		if err = dal.CreateFollow(ctx, &dal.Following{LeaderID: toUserID, FollowerID: userID}); err != nil {
			return util.ErrInternalError
		}
		if err = dal.UpdateDelta(ctx, toUserID, userID, 1); err != nil {
			return util.ErrInternalError
		}
	}
	return nil
}

// Unfollow 参数默认是合法的，即用户存在且操作合法，这个操作要4次左右的数据库查询，前端没有做重复请求拦截，这个有必要在外层优化
func Unfollow(ctx context.Context, req *user.UnfollowRequest) error {
	toUserID := req.LeaderId
	userID := req.FollowerId

	if toUserID == userID || toUserID == 0 || userID == 0 {
		return util.ErrIllegalArguments
	}
	follow, err := dal.QueryFollowUnscoped(ctx, toUserID, userID)
	if err != nil {
		return err
	}
	if follow.ID > 0 && !follow.Model.DeletedAt.Valid {
		if err = dal.DeleteFollow(ctx, follow); err != nil {
			return util.ErrInternalError
		}
		if err = dal.UpdateDelta(ctx, toUserID, userID, -1); err != nil {
			return util.ErrInternalError
		}
	}
	return nil
}

func FollowList(ctx context.Context, req *user.FollowListRequest) ([]*user.User, error) {
	userID := req.Req.UserId
	if userID == 0 {
		return nil, util.ErrIllegalArguments
	}
	repoUsers, err := dal.FollowList(ctx, userID)
	if err != nil {
		return nil, err
	}
	users := Users(repoUsers)
	for _, user := range users {
		user.IsFollow = true
	}
	return users, nil
}

func FollowerList(ctx context.Context, req *user.FollowerListRequest) ([]*user.User, error) {
	userID := req.Req.UserId
	followers, err := dal.FollowerList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return merge(ctx, followers, userID)
}
