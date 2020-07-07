package main

import (
	"net"
	"time"
	"log"
	"golang.org/x/sys/unix"
	"github.com/x2y2/poller-demo/poller"
	"github.com/x2y2/poller-demo/netpoll"
)

func start(ep *poller.Epoll){
	var buf = make([]byte,32)
	for {
		conns, err := ep.Wait()
		if err != nil{
			log.Println(err)
			continue
		}


		time.Sleep(time.Millisecond * 5)

		for _,conn := range conns{
			defer conn.Close()

			if poller.InEvents != 0 {
			     n :=  conn.Read(buf)
			     if n <= 0 || err != nil{
			        if err == unix.EAGAIN{return}
			        conn.Close()
			        break
			     }
			     log.Printf(string(buf))
			    }

			if poller.OutEvents != 0{
			err := conn.Write("welcome")
			if err != nil {
			    if err == unix.EAGAIN{return}
			    conn.Close()
			    break
				}
			}

			if poller.ErrEvents != 0{
				conn.Close()
				ep.Remove(conn)
			}
		}
	}
}


func main(){
	ln ,_ := net.Listen("tcp",":8000")
	log.Println("Server is listen at 8000")

	f ,_ := ln.(*net.TCPListener).File()
	fd := int(f.Fd())

	var ep *poller.Epoll
    ep ,_ = poller.NewEpoll()

    go start(ep)

    for{
    	nfd,sa,_ := unix.Accept(fd)
    	if err := unix.SetNonblock(nfd, true); err != nil {
            log.Println(err)
            return
        }

        conn := netpoll.NewTcpConn(nfd,sa)

        log.Printf("connect is coming from %s",conn.RemoteAddr)
        if err := ep.Add(*conn); err != nil{
        	log.Println(err)
        	conn.Close()
        }
    }
}