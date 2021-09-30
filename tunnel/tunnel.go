package tunnel

import (
	"io"
        "net"
        "log"
        "time"
        "github.com/zztczcx/go-dip/pool"
)

func HandleConn(source *net.TCPConn) {
	ip, err := pool.ChooseIp()
	if err != nil {
		source.Close()
                log.Printf("choose ip failed: %s", err.Error())
		return
	}

	dest, err := net.Dial("tcp", ip)
	if err != nil {
		source.Close()
		log.Printf("connect to %s failed: %s", ip, err.Error())
		return
	}

	source.SetKeepAlive(true)
	source.SetKeepAlivePeriod(time.Second * 60)

	go forward(source, dest.(*net.TCPConn))
	go forward(dest.(*net.TCPConn), source)
}

func forward(source *net.TCPConn, dest *net.TCPConn) {
	defer dest.CloseWrite()
	defer source.CloseRead()
	io.Copy(dest, source)
}
