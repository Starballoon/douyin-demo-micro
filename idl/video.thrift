namespace go video

include "user.thrift"

struct Video{
    1:i64 id
    2:user.User author
    3:string play_url
    4:string cover_url
    5:string title
    6:i64 favorite_count
    7:i64 comment_count
    8:bool is_favorite
}

struct MultiVideoResponse{
    1:list<Video> videos
    2:user.BaseResp resp
}

struct CreateVideoRequest{
    1:i64 author_id
    2:string video_filename
    3:string cover_filename
    4:string title
}

struct FeedRequest{
    1:optional i64 latest_time
    2:optional i64 user_id
}

struct FeedResponse{
    1:list<Video> video_list
    2:optional i64 next_time
    3:user.BaseResp resp
}

struct MGetVideoRequest{
    1:user.IdRequest req
}

struct MGetVideoResponse{
    1:MultiVideoResponse resp
}

struct PublishListRequest{
    1:user.IdRequest req
}

struct PublishListResponse{
    1:MultiVideoResponse resp
}

struct FavoriteListRequest{
    1:user.IdRequest req
}

struct FavoriteListResponse{
    1:MultiVideoResponse resp
}

struct FavoriteRequest{
    1:i64 video_id
    2:i64 user_id
}

struct FavoriteResponse{
    1:user.BaseResp resp
}

struct UnfavoriteRequest{
    1:i64 video_id
    2:i64 user_id
}

struct UnfavoriteResponse{
    1:user.BaseResp resp
}

struct CheckVideoRequest{
    1:i64 video_id
}

struct UpdateCommentCountRequest{
    1:i64 video_id
    2:i64 delta
}

service VideoService{
    user.BaseResp CreateVideo(1:CreateVideoRequest req)
    FeedResponse Feed(1:FeedRequest req)
    MGetVideoResponse MGetVideo(1:MGetVideoRequest red)
    PublishListResponse PublishList(1:PublishListRequest req)
    FavoriteListResponse FavoriteList(1:FavoriteListRequest req)
    FavoriteResponse Favorite(1:FavoriteRequest req)
    UnfavoriteResponse Unfavorite(1:UnfavoriteRequest req)
    user.BaseResp CheckVideo(1:CheckVideoRequest req)
    user.BaseResp UpdateCommentCount(1:UpdateCommentCountRequest req)
}