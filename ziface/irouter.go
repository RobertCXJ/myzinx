package ziface

// IRouter 路由抽象接口
// 路由里的数据都是IRequest
type IRouter interface {
	// PreHandle 处理conn之前的钩子方法Hook
	PreHandle(request IRequest)
	// Handle 处理conn业务的主方法Hook
	Handle(request IRequest)
	// PostHandle 处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
