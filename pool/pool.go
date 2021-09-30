package pool

import (
  "errors"
  "time"
  "encoding/json"
  "os"
  "strings"
  "fmt"
  "log"
)

type Ip struct {
  addr string
  weight float32
}

type Proxy struct {
  Ip string
  Port string
  Protocols []string
  Latency float32
}

var ips = make([]Ip, 0)

var queue = make(chan string)
var push = make(chan Ip)
var rotate = make(chan Ip)

// fetch ips from external
// test the ip's speed
// maintain a pool
func Run(){
  go fetchIp()
  go rotateIp()
  go appendIp()

  //waitGroup
}

func ChooseIp() (string, error){
  if len(ips) == 0 {
    return "", errors.New("ip pool is empty")
  }else{
    return <-queue, nil
  }
}

var rotateDone = make(chan bool)

func rotateIp() {
  for {
    if len(ips) == 0 {
      time.Sleep(10*time.Second)

    }else {
      <- rotateDone
      ip := ips[0]
      fmt.Println("ip is ", ip)
      queue <- ip.addr
      rotate <- ip
    }
  }
}
 
func fetchIp() error{
  fp, err := os.Open("./Free_Proxy_List.json")
  defer fp.Close()

  if err != nil {
      log.Printf("open json file failed: %s", err.Error())
      return err
  }

  var proxies []Proxy 

  dec := json.NewDecoder(fp)
  err = dec.Decode(&proxies)
  if err != nil {
      log.Printf("decode json file failed: %s", err.Error())
      return err
  }

  for _, proxy := range proxies {
    var sb strings.Builder
    sb.WriteString(proxy.Ip)
    sb.WriteString(":")
    sb.WriteString(proxy.Port)

    ip := Ip{addr: sb.String(), weight: proxy.Latency}
    push <- ip
  }

  rotateDone <- true


  for {
    time.Sleep(time.Second * 300)
    //check proxy speed
  }
}

func appendIp(){
  for {
    select {

      case ip := <- rotate:
        ips = append(ips[1:], ip)
        rotateDone <- true
      case ip := <- push:
        ips = append(ips, ip)
    } 
  }
}
