# douyin-demo 菜鸡互啄队

## 项目简介
主要使用KiteX、gin和GORM框架，在[官方demo](https://github.com/RaymondCode/simple-demo)的基础上参考[EasyNote项目](https://github.com/cloudwego/kitex-examples/tree/main/bizdemo/easy_note)将单体服务转换为微服务架构。根据项目接口和功能依赖关系，拆分为3个微服务和1个网关。项目中采用Docker部署MySQL数据库、OSS存储、Etcd注册中心和Redis。OSS资源由容器MinIO托管，Redis主要用于缓存用户操作。
## 接口描述

| 接口                              | 完成情况 | 入参                | 备注                          |
|---------------------------------|------|-------------------|-----------------------------|
| /static/                        | 100% |                   | 切换为OSS存储                    |
| /douyin/feed/                   | 100%  | token             |                             |
| /douyin/user/register           | 100%  | username/password |                             |
| /douyin/user/login              | 100%  | username/password |                             |
| /douyin/user/                   | 100%  | user_id/token     | user_id和token应该相符           | 
| /douyin/publish/action/         | 100%  | token             | 没有user_id                   |
| /douyin/publish/list/           | 100%  | user_id/token     |               |
| /douyin/favorite/action/        | 100%  | user_id/token     |                             |
| /douyin/favorite/list/          | 100%  | user_id/token     | 前端显示的数量不对，刷视频会拉当前视频作者 |               
| /douyin/comment/action/         | 100%  | user_id/token     | 前端只给了video_id，没有user_id     |
| /douyin/comment/list/           | 100%  | token             | 前端显示的评论总数不对，刷视频会拉当前视频作者     |
| /douyin/relation/action/        | 100%  | token             | 前端只给了to_user_id，没有user_id   |
| /douyin/relation/follow/list/   | 100%  | user_id/token     | 项目文档接口拼错了                  |
| /douyin/relation/follower/list/ | 100%  | user_id/token     |                             |
## 运行说明
本项目主要在Windows上使用WSL2和Docker Desktop相关虚拟机进行开发。

在项目根目录运行下列命令即可启动服务，但是由于容器启动较慢，建议等容器初始化完成再运行项目。
```shell
docker compose up -d
bash run.sh
```

项目IDL生成示例
```shell
kitex -module douyin-demo-micro -service user user.thrift
```

管理员权限运行将本机8080端口映射到WSL的8080端口
```shell
# 创建
# netsh interface portproxy add v4tov4 listenport=[win10端口] listenaddress=0.0.0.0 connectport=[虚拟机的端口] connectaddress=[虚拟机的ip]
netsh interface portproxy add v4tov4 listenport=8080 listenaddress=0.0.0.0 connectport=8080 connectaddress=172.24.171.176
# 删除
netsh interface portproxy delete v4tov4 listenport=8080 listenaddress=0.0.0.0
```

接口在线文档[抖音极简版-真菜鸡互啄队](https://www.apifox.cn/apidoc/project-1066782/api-22446795)

与配置有关的部分都硬编码放到util/config.go文件，未出现的public文件夹下是模拟用的数据文件

MySQL数据库定义了5张表，Video/User/Comment/User-User/Video-User，E-R图见PPT。因为数据比较简单，加了很多索引。

gorm里面的等价操作
```go
// 对于单主键查询，First比Find略快，但是不显著，需要进一步测试
DB.Find(&User, 1)
DB.Where("id=?", 1).Find(&User)

DB.Find(&User, []int{1,2})
DB.Where("id in ?", []int{1,2}).Find(&User)

DB.Delete(&User, 1)
DB.Where("id=?", 1).Delete(&User)
```

service部分的Video增加了Title字段，而且前端确实有这个字段，是视频的标题。

本项目开发期间魔改了"github.com/appleboy/gin-jwt/v2"库的auth_jwt.go里面的ParseToken，增加了一个从request body的form获取token的方式，为了兼容所有接口。
目前相关PR已合并，go.mod 需要手动设置依赖为 github.com/appleboy/gin-jwt/v2 v2.8.1-0.20220605135842-8f9474155532 及以上。

### 功能说明

* 静态资源模式下(util.STATIC=true)视频上传后会保存到本地 public 目录中，访问时用ip:8080/static/videos/video_name，封面在static/covers/
* OSS存储模式下(util.STATIC=false)，视频和封面会托管到MinIO，需要在本地生成一次封面，需要使用```docker compose up```启动若干容器

### 测试数据

测试数据写在若干fakedata.go