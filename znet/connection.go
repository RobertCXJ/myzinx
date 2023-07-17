package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool

	// 通知当前连接停止的 channel
	ExitChan chan bool

	// 该连接处理的方法Router
	Router ziface.IRouter
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// StartReader 处理连接读数据的 goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, ", Reader is exit, remote addr is ", c.RemoteAddr().String()) // 2
	defer c.Stop()                                                                                         // 1

	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error ", err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 从路由中，找到注册绑定的connection对应的router调用
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnId = ", c.ConnID)
	// 启动从当前连接读数据的 goroutine
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()    // 关闭socket连接
	close(c.ExitChan) // 回收资源
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
