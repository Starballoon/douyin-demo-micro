# douyin-demo 菜鸡互啄队

## 修改的主要内容

| 接口                              | 完成情况 | 入参                | 备注                          |
|---------------------------------|------|-------------------|-----------------------------|
| /static/                        | 100% |                   | 切换为OSS存储                    |
| /douyin/feed/                   | 50%  | token             |                             |
| /douyin/user/register           | 70%  | username/password |                             |
| /douyin/user/login              | 60%  | username/password |                             |
| /douyin/user/                   | 70%  | user_id/token     | user_id和token应该相符           | 
| /douyin/publish/action/         | 80%  | token             | 没有user_id                   |
| /douyin/publish/list/           | 60%  | user_id/token     | 可以拉取其他人的发布列表吗？              |
| /douyin/favorite/action/        | 30%  | user_id/token     |                             |
| /douyin/favorite/list/          | 30%  | user_id/token     | 数量无限制，前端显示的数量不对，刷视频会拉当前视频作者 |               
| /douyin/comment/action/         | 30%  | user_id/token     | 前端只给了video_id，没有user_id     |
| /douyin/comment/list/           | 30%  | token             | 前端显示的评论总数不对，刷视频会拉当前视频作者     |
| /douyin/relation/action/        | 30%  | token             | 前端只给了to_user_id，没有user_id   |
| /douyin/relation/follow/list/   | 30%  | user_id/token     | 项目文档接口拼错了吧                  |
| /douyin/relation/follower/list/ | 30%  | user_id/token     |                             |

docker文件还是有些问题
要运行的话可以在WSL之类的Linux的项目根目录运行下列命令
```shell
docker compose up -d
bash run.sh
```

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

增加了docker操作，计划改成微服务架构之后需要托管到容器中去。

在项目根目录下运行，需要安装docker，编译服务镜像，注意config.go中的一些配置，目前由于客户端不支持docker而无法正常使用。
```shell
docker build -t douyin-demo-server .
```

有些接口会同时提供user_id/token，如何两者校验不一致，是否要考虑缓存非法请求的ip、设备号之类的内容在授权时就禁用。

项目没有考虑重复点击的问题，可能需要根据服务ID+用户请求生成一个唯一标识，做一个LRU。或者每个service层的端口单独缓存正在执行的服务？

配置部分都放到util模块了，public文件夹下是数据文件

controller/common.go -> service/common.go 避免循环引用。增加了repository用于简化数据库部分。

MySQL数据库定义了5张表，Video/User/Comment/User-User/Video-User，E-R图见PPT。因为数据比较简单，加了很多索引。

索引加了之后多了很多问题，因为gorm默认是软删除，如果要覆盖这部分逻辑还是要手动改一下。

gorm里有个坑，当使用了Model方法，且该对象主键有值，该值会被用于构建条件。

本项目涉及的级联查询不知道是否可以采用gorm的preload机制做缓存加速[预加载](https://gorm.io/zh_CN/docs/preload.html)，或者改用复杂查询。

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

魔改了"github.com/appleboy/gin-jwt/v2"库的auth_jwt.go里面的ParseToken，增加了一个从request body的form获取token的方式，为了兼容所有接口。不确定其他参数行不行。
提的PR已合并，go.mod 设置依赖为 github.com/appleboy/gin-jwt/v2 v2.8.1-0.20220605135842-8f9474155532 及以上即可。

### 功能说明

* 静态资源模式下(util.STATIC=true)视频上传后会保存到本地 public 目录中，访问时用ip:8080/static/videos/video_name，封面在static/covers/
* OSS存储模式下(util.STATIC=false)，视频和封面会托管到MinIO，需要在本地生成一次封面，需要使用```docker compose up```启动若干容器

### 测试数据

测试数据写在若干fakedata.go