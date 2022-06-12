package rpc

// InitRPC init rpc clients
func InitRPC() {
	initUserRPC()
	initVideoRPC()
	initCommentRPC()
	initMinIO()
}
