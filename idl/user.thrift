namespace go user

struct User{
    1:i64 id
    2:string name
    3:i64 follow_count
    4:i64 follower_count
    5:bool is_follow
}
struct NameRequest{
    1:string username
    2:string password
}

struct IdRequest {
    1:i64 user_id
}

struct BaseResp{
    1:i64 status_code
    2:string status_message
    3:i64 service_time
}

struct UserResponse{
    1:User user
    2:BaseResp resp
}

struct MultiUserResponse{
    1:list<User> users
    2:BaseResp resp
}

struct CreateUserRequest {
    1:NameRequest req
}

struct CreateUserResponse{
    1:i64 user_id
    2:BaseResp resp
}

struct CheckUserRequest{
    1:NameRequest req
}

struct CheckUserResponse{
    1:UserResponse resp
}

struct FindUserRequest {
    1:IdRequest req
}

struct FindUserResponse{
    1:UserResponse resp
}

struct MGetUserRequest{
    1:list<i64> user_ids
    2:optional i64 user_id
}

struct MGetUserResponse{
    1:MultiUserResponse resp
}

struct FollowRequest{
    1:i64 leader_id
    2:i64 follower_id
}

struct FollowResponse{
    1:BaseResp resp
}

struct UnfollowRequest{
    1:i64 leader_id
    2:i64 follower_id
}

struct UnfollowResponse{
    1:BaseResp resp
}

struct FollowListRequest{
    1:IdRequest req
}

struct FollowListResponse{
    1:MultiUserResponse resp
}

struct FollowerListRequest{
    1:IdRequest req
}

struct FollowerListResponse{
    1:MultiUserResponse resp
}

service UserService{
    CreateUserResponse CreateUser(1:CreateUserRequest req)
    CheckUserResponse CheckUser(1:CheckUserRequest req)
    FindUserResponse FindUser(1:FindUserRequest req)
    MGetUserResponse MGetUser(1:MGetUserRequest req)
    FollowResponse Follow(1:FollowRequest req)
    UnfollowResponse Unfollow(1:UnfollowRequest req)
    FollowListResponse FollowList(1:FollowListRequest req)
    FollowerListResponse FollowerList(1:FollowerListRequest req)
}
