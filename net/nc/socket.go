/**
 * Copyright 2015 @ to2.net.
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
	"sync"
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
	AuthFunc func() (from int64, err error)

	// 客户端
	Client struct {
		// 连接来源,如:来源与某个商户
		Source int64
		// 用户编号,如:商户下面的客户编号
		User int64
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
		output       bool               // 开启输出,默认开启
		ReadDeadLine time.Duration      //超时断开时间
		clients      map[string]*Client //客户端身份,断开时会删除
		userAddrs    map[int64]string   //用户的IP信息,以同时下发到多个客户端
		handlers     map[string]CmdFunc
		jobs         []Job
		tr           TcpReceiver
		mux          sync.Mutex
	}
)

func NewSocketServer(r TcpReceiver) *SocketServer {
	return &SocketServer{
		output:    true,
		clients:   make(map[string]*Client),
		userAddrs: make(map[int64]string),
		jobs:      make([]Job, 0),
		tr:        r,
	}
}

func (s *SocketServer) OutputOff() {
	s.output = false
}

// print
func (s *SocketServer) Printf(format string, args ...interface{}) {
	if s.output {
		log.Printf(format, args...)
	}
}

// Register job running before server start!
func (s *SocketServer) RegisterJob(job Job) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, v := range s.jobs {
		if reflect.ValueOf(v) == reflect.ValueOf(job) {
			return errors.New("Can't repeat register job!")
		}
	}
	s.jobs = append(s.jobs, job)
	return nil
}

// Unregister job
func (s *SocketServer) UnregisterJob(job Job) {
	for i, v := range s.jobs {
		if reflect.ValueOf(v) == reflect.ValueOf(job) {
			// compare func ptr
			s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
			break
		}
	}
}

// 根据IP获取客户端信息,如果没有,返回false
func (s *SocketServer) GetCli(conn net.Conn) (*Client, bool) {
	c, b := s.clients[conn.RemoteAddr().String()]
	return c, b
}

// 获取用户的客户端连接, 一个用户对应一个或多个连接
func (s *SocketServer) GetConnections(userId int64) []net.Conn {
	arr := strings.Split(s.userAddrs[userId], "$")
	var connList []net.Conn = make([]net.Conn, 0)
	for _, v := range arr {
		if i, ok := s.clients[v]; ok && i.Conn != nil {
			connList = append(connList, i.Conn)
		}
	}
	return connList
}

// 验证消息
func (s *SocketServer) Auth(conn net.Conn, f AuthFunc) error {
	var src int64
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
		s.clients[addr] = cli
		s.Printf("[ CLIENT] - Auth Success! source %d(%s)", src, addr)
	}
	return err
}

// 用户验证
func (s *SocketServer) UAuth(conn net.Conn, f AuthFunc) error {
	var uid int64
	var err error
	var cli *Client
	if f != nil {
		uid, err = f()
	}
	if err == nil {
		if cli, _ = s.GetCli(conn); cli != nil {
			cli.User = uid
			cli.LatestConnectTime = time.Now()
			//设置用户连接的客户端,一个用户可能连接多个终端
			if v, ok := s.userAddrs[uid]; ok {
				s.userAddrs[uid] = v + "$" + cli.Addr.String()
			} else {
				s.userAddrs[uid] = cli.Addr.String()
			}
		}
	}
	return err
}

// 不需要验证
func (s *SocketServer) NoAuth(conn net.Conn) error {
	return s.Auth(conn, nil)
}

// 存储客户端信息,通常需要通过某种权限后才允许访问SOCKET服务
func (s *SocketServer) SetCli(conn net.Conn, c *Client) {
	s.clients[conn.RemoteAddr().String()] = c
}

//func (s *SocketServer) setUserAddrs(id int64, addr string) {
//	s.userAddrs[id] = addr
//}

func (s *SocketServer) runJobs() {
	for _, v := range s.jobs {
		go v(s)
	}
}

func (s *SocketServer) Listen(addr string) {
	s.runJobs() //running job

	serveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	listen, err := net.ListenTCP("tcp", serveAddr)
	for {
		if conn, err := listen.AcceptTCP(); err == nil {
			s.Printf("[ CLIENT][ CONNECT] - New client %s ; actived clients : %d",
				conn.RemoteAddr().String(), len(s.clients)+1)
			go s.receiveTcpConn(conn, s.tr)
		}
	}
}

// Receive client connection
func (s *SocketServer) receiveTcpConn(conn *net.TCPConn, rc TcpReceiver) {
	const delim byte = '\n'
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadBytes(delim)

		if err != nil {
			addr := conn.RemoteAddr().String()
			//断开连接,清理数据
			//一个用户可能通过不同的地址连接
			if v, ok := s.clients[addr]; ok {
				uid := v.User
				delete(s.clients, addr)
				addr2 := s.userAddrs[uid] //获取用户所有的客户端地址
				if strings.Index(addr2, "$") == -1 {
					//清除用户终端地址
					delete(s.userAddrs, uid)
				} else {
					addr2 = strings.Replace(addr2, addr, "", 1)
					s.userAddrs[uid] = strings.Replace(addr2, "$$", "$", -1)
				}
			}
			s.Printf("[ CLIENT][ DISCONN] - Client %s disconnect, actived clients : %d",
				addr, len(s.clients))
			break
		}

		// custom handle client request
		if d, err := rc(conn, line[:len(line)-1]); err != nil {
			// remove '\n'
			conn.Write([]byte("error$" + err.Error()))
		} else if d != nil {
			conn.Write(d)
		}

		conn.Write([]byte{delim})
		if s.ReadDeadLine > 0 {
			// set connect time out
			conn.SetReadDeadline(time.Now().Add(s.ReadDeadLine))
		}
	}
}
