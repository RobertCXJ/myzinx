package ziface

// IRequest 接口：
// 实际上是把 客户端请求的连接信息 和 请求的数据 包装到了一个Request中

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgID 得到当前请求消息 ID
	GetMsgID() uint32
}
