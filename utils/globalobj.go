package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// GlobalObj 存储一切有关Zinx框架的全局参数，供其他模块使用
// 一些参数通过 zinx.json 由用户进行配置
type GlobalObj struct {
	// Server
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	// Zinx
	Version        string // 当前Zinx的版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前Zinx框架数据包的最大值
}

// GlobalObject 定义一个全局的对外对象
var GlobalObject *GlobalObj

// 初始化当前的 GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.5",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	// 尝试从 conf/zinx.json 去加载用户自定义的参数
	//GlobalObject.Reload()
}

// Reload 从 zinx.json 加载自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
