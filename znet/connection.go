package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool

	// 通知当前连接停止的 channel(由 Reader 告知 Writer 退出)
	ExitChan chan bool

	// 无缓冲的管道，用于读、写Goroutine之间的通信
	msgChan chan []byte

	// 消息的管理 MsgID 和对应处理业务 API 关系
	MsgHandler ziface.IMsgHandle
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

// StartReader 处理连接读数据的 goroutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID = ", c.ConnID, ", [Reader is exit], remote addr is ", c.RemoteAddr().String()) // 2
	defer c.Stop()                                                                                           // 1

	for {
		// 读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf error ", err)
		//	continue
		//}
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的 Msg Head 二进制流 8 个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}
		// 拆包，得到 msgID 和 msgDataLen 放到 mgs 消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		// 根据 dataLen 再次读取 Data，放在 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池，将消息发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的connection对应的router调用
			// 根据绑定好的MsgID找到处理对应API业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写消息Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")
	// 不断的阻塞等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// Reader 已经退出，此时 Writer 也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnId = ", c.ConnID)
	// 启动从当前连接读数据的 goroutine
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭 socket 连接
	c.Conn.Close()
	// 告知 Writer 关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// SendMsg 将要发送给客户端的数据，先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	// MsgDataLen/MsgId/Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg is = ", msgId)
		return errors.New("pack err msg")
	}
	// 读写分离：将数据发送给消息管道
	c.msgChan <- binaryMsg
	return nil
}
