package poller

import (
	"log"
	"sync"
	"golang.org/x/sys/unix"
    "github.com/x2y2/poller-demo/netpoll"
)

const (
    ErrEvents = unix.EPOLLERR | unix.EPOLLHUP | unix.EPOLLRDHUP
    OutEvents = ErrEvents | unix.EPOLLOUT
    InEvents = ErrEvents | unix.EPOLLIN | unix.EPOLLPRI
)

type Epoll struct{
	fd 		int
	lock 	*sync.RWMutex
	conns  	map[int]netpoll.Conn
}

func NewEpoll()(*Epoll, error){
    fd,err := unix.EpollCreate1(0)
    if err != nil{
        return nil,err
    }

    return &Epoll{
        fd:     fd,
        lock:   &sync.RWMutex{},
        conns: 	make(map[int]netpoll.Conn),
    },nil
}

func (e *Epoll) Add(conn netpoll.Conn) error{
    fd := conn.Fd
    events := &unix.EpollEvent{Events: unix.POLLIN|unix.POLLHUP, Fd: int32(fd)}
    err := unix.EpollCtl(e.fd,unix.EPOLL_CTL_ADD,fd,events)
    if err != nil{
        return err
    }

    e.lock.Lock()
    defer e.lock.Unlock()
    e.conns[fd] = conn
    return nil
}

func (e *Epoll) Remove(conn netpoll.Conn) error{
    fd := conn.Fd
    err := unix.EpollCtl(e.fd,unix.EPOLL_CTL_DEL,fd,nil)
    if err != nil{
        return err
    }
    e.lock.Lock()
    defer e.lock.Unlock()
    delete(e.conns,fd)
    return nil
}

func (e *Epoll) Wait()([]netpoll.Conn,error){
    events := make([]unix.EpollEvent,100)
    n, err := unix.EpollWait(e.fd,events,-1)
    if err != nil{
        log.Println(err)
        return  nil,err
    }

    e.lock.RLock()
    defer e.lock.RUnlock()
    var conns []netpoll.Conn
    for i := 0;i < n; i++{
        c := e.conns[int(events[i].Fd)]
        conns = append(conns,c)
    }
    return conns,nil
}