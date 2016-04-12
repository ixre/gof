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
	"errors"
	"log"
	"net"
	"reflect"
	"strings"
	"time"
)

type (
	// TCP接收者
	TcpReceiver func(conn net.Conn, read []byte) ([]byte, error)

	// 服务端运行的Job
	Job func(*SocketServer)

	// 命令函数,@plan:客户端发送的报文
	CmdFunc func(c *Client, plan string) ([]byte, error)

	// 验证函数,返回编号及错误
	AuthFunc func() (from int, err error)

	// 客户端
	Client struct {
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

	// Socket服务器, 一个IP地址对应一个连接,一个用户对应一个或多个IP和连接.
	// 多个IP用"$"连接
	SocketServer struct {
		// 开启输出,默认开启
		_output      bool
		ReadDeadLine time.Duration      //超时断开时间
		_clients     map[string]*Client //客户端身份,断开时会删除
		_userAddrs   map[int]string     //用户的IP信息,以同时下发到多个客户端
		_handlers    map[string]CmdFunc
		_jobs        []Job
		_r           TcpReceiver
	}
)

func NewSocketServer(r TcpReceiver) *SocketServer {
	return &SocketServer{
		_output:    true,
		_clients:   make(map[string]*Client),
		_userAddrs: make(map[int]string),
		_jobs:      make([]Job, 0),
		_r:         r,
	}
}

func (this *SocketServer) OutputOff() {
	this._output = false
}

// print
func (this *SocketServer) Printf(format string, args ...interface{}) {
	if this._output {
		log.Printf(format, args...)
	}
}

// Register job running before server start!
func (this *SocketServer) RegisterJob(job Job) error {
	for _, v := range this._jobs {
		if reflect.ValueOf(v) == reflect.ValueOf(job) { // compare func ptr
			return errors.New("Can't repeat register job!")
		}
	}
	this._jobs = append(this._jobs, job)
	return nil
}

// Unregister job
func (this *SocketServer) UnregisterJob(job Job) {
	for i, v := range this._jobs {
		if reflect.ValueOf(v) == reflect.ValueOf(job) { // compare func ptr
			this._jobs = append(this._jobs[:i], this._jobs[i+1:]...)
			break
		}
	}
}

// 根据IP获取客户端信息,如果没有,返回false
func (this *SocketServer) GetCli(conn net.Conn) (*Client, bool) {
	c, b := this._clients[conn.RemoteAddr().String()]
	return c, b
}

// 获取用户的客户端连接, 一个用户对应一个或多个连接
func (this *SocketServer) GetConnections(userId int) []net.Conn {
	arr := strings.Split(this._userAddrs[userId], "$")
	var connList []net.Conn = make([]net.Conn, 0)
	for _, v := range arr {
		if i, ok := this._clients[v]; ok && i.Conn != nil {
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
		this._clients[addr] = cli
		this.Printf("[ CLIENT] - Auth Success! source %d(%s)", src, addr)
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
			if v, ok := this._userAddrs[uid]; ok {
				this._userAddrs[uid] = v + "$" + cli.Addr.String()
			} else {
				this._userAddrs[uid] = cli.Addr.String()
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
	this._clients[conn.RemoteAddr().String()] = c
}

func (this *SocketServer) setUserAddrs(id int, addr string) {
	this._userAddrs[id] = addr
}

func (this *SocketServer) runJobs() {
	for _, v := range this._jobs {
		go v(this)
	}
}

func (this *SocketServer) Listen(addr string) {
	this.runJobs() //running job

	serveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	listen, err := net.ListenTCP("tcp", serveAddr)
	for {
		if conn, err := listen.AcceptTCP(); err == nil {
			this.Printf("[ CLIENT][ CONNECT] - New client %s ; actived clients : %d",
				conn.RemoteAddr().String(), len(this._clients)+1)
			go this.receiveTcpConn(conn, this._r)
		}
	}
}

// Receive client connection
func (this *SocketServer) receiveTcpConn(conn *net.TCPConn, rc TcpReceiver) {
	const delim byte = '\n'
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadBytes(delim)

		if err != nil {
			addr := conn.RemoteAddr().String()
			//断开连接,清理数据
			//一个用户可能通过不同的地址连接
			if v, ok := this._clients[addr]; ok {
				uid := v.User
				delete(this._clients, addr)
				addr2 := this._userAddrs[uid]        //获取用户所有的客户端地址
				if strings.Index(addr2, "$") == -1 { //清除用户终端地址
					delete(this._userAddrs, uid)
				} else {
					addr2 = strings.Replace(addr2, addr, "", 1)
					this._userAddrs[uid] = strings.Replace(addr2, "$$", "$", -1)
				}
			}
			this.Printf("[ CLIENT][ DISCONN] - Client %s disconnect, actived clients : %d",
				addr, len(this._clients))
			break
		}

		// custom handle client request
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
