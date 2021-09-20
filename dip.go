package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
        "github.com/zztczcx/go-dip/pool"
        "github.com/zztczcx/go-dip/tunnel"
)


func init() {
        rand.Seed(time.Now().Unix())
}

const SIG_RELOAD = syscall.Signal(35)
const SIG_STATUS = syscall.Signal(36)

func status() {
	log.Printf("num goroutines: %d", runtime.NumGoroutine())
	buf := make([]byte, 32768)
	runtime.Stack(buf, true)
	log.Printf("!!!!!stack!!!!!:%s", buf)
}

func monitor(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, SIG_RELOAD, SIG_STATUS, syscall.SIGTERM, syscall.SIGHUP)

	for sig := range c {
		switch sig {
		case SIG_RELOAD:
			//reload()
		case SIG_STATUS:
			status()
		default:
			log.Printf("catch siginal: %v, ignored", sig)
		}
	}

}

func main() {
	var laddr string
	flag.StringVar(&laddr, "listen", ":1248", "local listen port")
	flag.Parse()

	go monitor()
        go pool.Run()
         
	// run
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Printf("build listener failed:%s", err.Error())
		return
	}
	defer ln.Close()


	// run loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept failed:%s", err.Error())
			if opErr, ok := err.(*net.OpError); ok {
				if !opErr.Temporary() {
					break
				}
			}
			continue
		}
		go tunnel.HandleConn(conn.(*net.TCPConn))
	}
}
