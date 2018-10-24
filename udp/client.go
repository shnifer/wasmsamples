package main

import (
	"net"
	"time"
	"strconv"
)

func main(){
	addr,err:=net.ResolveUDPAddr("udp","localhost:8001")
	if err!=nil{
		panic(err)
	}
	conn,err:=net.DialUDP("udp",nil,addr)
	if err!=nil{
		panic(err)
	}
	for i:=0; i<10;i++{
		conn.Write([]byte(strconv.Itoa(i)))
		time.Sleep(time.Second)
	}
}
