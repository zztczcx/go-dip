package main


func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var laddr string
	flag.StringVar(&laddr, "listen", ":1248", "local listen port")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Println("config file is missed.")
		return
	}

	options.config = args[0]
	if err := reloadConfig(); err != nil {
		log.Printf("load config failed:%v", err)
		return
	}

	// run
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Printf("build listener failed:%s", err.Error())
		return
	}
	defer ln.Close()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, SIG_RELOAD, SIG_STATUS, syscall.SIGTERM, syscall.SIGHUP)

		for sig := range c {
			switch sig {
			case SIG_RELOAD:
				reload()
			case SIG_STATUS:
				status()
			default:
				log.Printf("catch siginal: %v, ignored", sig)
			}
		}
	}()

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
		go handleConn(conn.(*net.TCPConn))
	}
}
