/**
 * Copyright 2015 @ z3q.net.
 * name : authed_socket
 * author : jarryliu
 * date : 2015-12-28 15:32
 * description : 根据原型改编,原型参见:tcpserve.go.prototype
 * history :
 */
package nc

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

type (
	// TCP接收者
	TcpReceiver func(conn net.Conn, read []byte) ([]byte, error)
	// 命令函数,@plan:客户端发送的报文
	CmdFunc func(c *Client, plan string) ([]byte, error)
	// 验证函数,返回编号及错误
	AuthFunc func() (from int, err error)
	Client   struct {
		// 连接来源,如:来源与某个商户
		Source int
		// 用户编号,如:商户下面的客户编号
		User int
		// 客户端连接地址
		Addr net.Addr
		// 连接
		Conn net.Conn
		// 连接时间
		CreateTime time.Time
		// 最近连接时间
		LatestConnectTime time.Time
	}

	// Socket服务器
	SocketServer struct {
		// 开启输出,默认开启
		output       bool
		ReadDeadLine time.Duration      //超时断开时间
		clients      map[string]*Client //客户端身份
		userAddrs    map[int]string     //用户的IP信息,以同时下发到多个客户端
		handlers     map[string]CmdFunc
	}
)

func NewSocketServer() *SocketServer {
	return &SocketServer{
		output:    true,
		clients:   make(map[string]*Client),
		userAddrs: make(map[int]string),
	}
}

func (this *SocketServer) OutputOff() {
	this.output = false
}

// print logtime.Minute * 5
func (this *SocketServer) Print(format string, args ...interface{}) {
	if this.output {
		log.Printf(format, args...)
	}
}

// 获取客户端信息,如果没有,返回false
func (this *SocketServer) GetCli(conn net.Conn) (*Client, bool) {
	c, b := this.clients[conn.RemoteAddr().String()]
	return c, b
}

// 获取用户的客户端连接
func (this *SocketServer) GetConn(userId int) []net.Conn {
	arr := strings.Split(this.userAddrs[userId], "$")
	var connList []net.Conn = make([]net.Conn, 0)
	for _, v := range arr {
		if i, ok := this.clients[v]; ok && i.Conn != nil {
			connList = append(connList, i.Conn)
		}
	}
	return connList
}

// 验证消息
func (this *SocketServer) Auth(conn net.Conn, f AuthFunc) error {
	var src int
	var err error
	if f != nil {
		src, err = f()
	}
	if err == nil {
		now := time.Now()
		cli := &Client{
			Source:            src,
			User:              0,
			Addr:              conn.RemoteAddr(),
			Conn:              conn,
			CreateTime:        now,
			LatestConnectTime: now,
		}
		addr := cli.Addr.String()
		this.clients[addr] = cli
		this.Print(false, "[ CLIENT] - Auth Success! source %d(%s)", src, addr)
	}
	return err
}

func (this *SocketServer) UAuth(conn net.Conn, f AuthFunc) error {
	var uid int
	var err error
	var cli *Client
	if f != nil {
		uid, err = f()
	}
	if err == nil {
		if cli, _ = this.GetCli(conn); cli != nil {
			cli.User = uid
			cli.LatestConnectTime = time.Now()
			//设置用户连接的客户端,一个用户可能连接多个终端
			if v, ok := this.userAddrs[uid]; ok {
				this.userAddrs[uid] = v + "$" + cli.Addr.String()
			} else {
				this.userAddrs[uid] = cli.Addr.String()
			}
		}
	}
	return err
}

// 不需要验证
func (this *SocketServer) NoAuth(conn net.Conn) error {
	return this.Auth(conn, nil)
}

// 存储客户端信息,通常需要通过某种权限后才允许访问SOCKET服务
func (this *SocketServer) SetCli(conn net.Conn, c *Client) {
	this.clients[conn.RemoteAddr().String()] = c
}

func (this *SocketServer) setUserAddrs(id int, addr string) {
	this.userAddrs[id] = addr
}

func (this *SocketServer) Listen(addr string, rc TcpReceiver) {
	serveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp", serveAddr)
	for {
		if conn, err := listen.AcceptTCP(); err == nil {
			this.Print(false, "[ CLIENT][ CONNECT] - New client %s ; actived clients : %d",
				conn.RemoteAddr().String(), len(this.clients)+1)
			go this.receiveTcpConn(conn, rc)
		}
	}
}

func (this *SocketServer) receiveTcpConn(conn *net.TCPConn, rc TcpReceiver) {
	const delim byte = '\n'
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadBytes(delim)

		if err != nil {
			addr := conn.RemoteAddr().String()
			//断开连接,清理数据
			//一个用户可能通过不同的地址连接
			if v, ok := this.clients[addr]; ok {
				uid := v.User
				delete(this.clients, addr)
				addr2 := this.userAddrs[uid]         //获取用户所有的客户端地址
				if strings.Index(addr2, "$") == -1 { //清除用户终端地址
					delete(this.userAddrs, uid)
				} else {
					addr2 = strings.Replace(addr2, addr, "", 1)
					this.userAddrs[uid] = strings.Replace(addr2, "$$", "$", -1)
				}
			}
			this.Print(false, "[ CLIENT][ DISCONN] - Client %s disconnect, activeed clients : %d",
				addr, len(this.clients))
			break
		}

		if d, err := rc(conn, line[:len(line)-1]); err != nil { // remove '\n'
			conn.Write([]byte("error$" + err.Error()))
		} else if d != nil {
			conn.Write(d)
		}

		conn.Write([]byte{delim})
		if this.ReadDeadLine > 0 { // set connect time out
			conn.SetReadDeadline(time.Now().Add(this.ReadDeadLine))
		}
	}
}
