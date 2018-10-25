package main

import (
	"net"
	"time"
	"strconv"
	"log"
)

const msgSize = 1000

func listenLoop(conn *net.UDPConn){
	buf := make([]byte, msgSize)
	for{
		n, addr, err:= conn.ReadFromUDP(buf)
		if err!=nil {
			log.Println("ReadFromUDP err: ",err)
			continue
		}
		log.Println("recieved",n,"bytes from", addr,": ",string(buf[:n]))
	}
}

func main(){
	addr,err:=net.ResolveUDPAddr("udp","localhost:8001")
	if err!=nil{
		panic(err)
	}
	conn,err:=net.DialUDP("udp",nil,addr)
	if err!=nil{
		panic(err)
	}
	go listenLoop(conn)
	reallyBigBuff:=make([]byte, 10000)
	reallyBigBuff[0]='x'
	reallyBigBuff[9999]='z'
	conn.Write(reallyBigBuff)
	for i:=0; i<10;i++{
		conn.Write([]byte(strconv.Itoa(i)))
		time.Sleep(time.Second)
	}

}
