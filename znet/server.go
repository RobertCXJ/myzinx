package znet

import (
	"fmt"
	net2 "net"
	"zinx/ziface"
)

// Server IServer的接口实现, 定义一个Server的服务器模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int

	// 当前的Server添加一个router，server注册的连接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[START] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1 获取一个TCP的Addr
		addr, err := net2.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址
		listener, err := net2.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		// 监听成功
		fmt.Println("start Zinx server, ", s.Name, "succ, Listening...")
		var cid uint32 = 0
		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			// 3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的

			// 将处理新连接的业务方法和conn进行绑定，得到我们定义的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//TODO 将一些服务器的资源、状态、或者 已经开辟的连接信息进行停止或回收
}

func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	//TODO 做一些启动服务后的额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
