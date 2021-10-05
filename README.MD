## Go-Dynamic-IP

### Target

- It will port forward the connection to another proxy ip. 

- This service will auto update and test the available proxy ip pool.

### How to run it

  - through docker 
    * docker run zztczcx/go-dip 
  - build from Dockerfile
    *  docker build -t go-dip .
    *  docker run --publish 1248:1248 --name go-dip --rm go-dip
  - build from source
    * go build 
    * ./go-dip 
  
  By default the open port is 1248

  nc localhost 1248 ----> the output is tunneled to another ip
