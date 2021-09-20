package tunnel

import (
	"io"
        "net"
        "log"
        "time"
        "github.com/zztczcx/go-dip/pool"
)

func Forward(source *net.TCPConn, dest *net.TCPConn) {
	defer dest.CloseWrite()
	defer source.CloseRead()
	io.Copy(dest, source)
}
func HandleConn(source *net.TCPConn) {
	host := pool.Host()
	if host == nil {
		source.Close()
		log.Println("choose host failed")
		return
	}

	dest, err := net.Dial("tcp", host)
	if err != nil {
		source.Close()
		log.Printf("connect to %s failed: %s", host, err.Error())
		return
	}

	source.SetKeepAlive(true)
	source.SetKeepAlivePeriod(time.Second * 60)

	go Forward(source, dest.(*net.TCPConn))
	go Forward(dest.(*net.TCPConn), source)
}
