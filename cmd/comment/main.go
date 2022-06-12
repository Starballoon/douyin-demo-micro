package main

import (
	"douyin-demo-micro/cmd/comment/dal"
	"douyin-demo-micro/cmd/comment/rpc"
	comment "douyin-demo-micro/kitex_gen/comment/commentservice"
	"douyin-demo-micro/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"net"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{util.EtcdAddress})
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1"+util.CommentServicePort)
	if err != nil {
		panic(err)
	}

	util.InitJaeger(util.CommentService)
	rpc.InitRPC()
	err = dal.InitDB()
	if err != nil {
		panic(err)
	}

	svr := comment.NewServer(new(CommentServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: util.CommentService}), // server name
		server.WithMiddleware(util.CommonMiddleware),                                             // middleware
		server.WithMiddleware(util.ServerMiddleware),
		server.WithServiceAddr(addr),                                       // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                    // tracer
		server.WithBoundHandler(util.NewCpuLimitHandler()),                 // BoundHandler
		server.WithRegistry(r),                                             // registry)
	)
	err = svr.Run()

	if err != nil {
		klog.Fatal(err)
	}
}
