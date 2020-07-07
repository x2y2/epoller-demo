package main 

import (
    // "fmt"
    "net"
    "flag"
    "log"
    "time"
)

var (
    ip = flag.String("ip","192.168.20.23","server ip")
    connections = flag.Int("conn",1000,"number of tcp connections")
)

func main(){
    flag.Parse()

    addr := *ip + ":8000"
    log.Printf("connect to %s",addr)

    var conns []net.Conn

    for i := 0;i < *connections; i++{
        c,err := net.DialTimeout("tcp",addr, 5 * time.Second)
        if err != nil {
            log.Println("failed to connecte",i,err)
            i--
            continue
        }
        conns = append(conns,c)
        time.Sleep(time.Millisecond)
    }
    defer func(){
        for _, c := range conns{
            c.Close()
        }
    }()

    log.Printf("complete initiate %d connections",len(conns))
    tts := time.Second
    if *connections > 100{
        tts = time.Millisecond * 5
    }
    var buf = make([]byte,32)

    for i := 0; i < len(conns); i++{
        time.Sleep(tts)
        conn := conns[i]
        conn.Write([]byte("hello world"))
        conn.Read(buf)
        log.Println(string(buf))
    }
}