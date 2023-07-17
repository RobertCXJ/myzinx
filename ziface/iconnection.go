package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接 - 让当前的连接准备开始工作
	Start()
	// Stop 停止连接 - 结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接绑定的 socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接的连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的 TCP状态 IP Port
	RemoteAddr() net.Addr
	// Send 发送数据给远程客户端
	Send(data []byte) error
}

// HandleFunc 定义一个统一处理连接业务的接口
// 参数1: 原生socket连接
// 参数2: 客户端请求的数据
// 参数3: 客户端请求数据长度
type HandleFunc func(*net.TCPConn, []byte, int) error
