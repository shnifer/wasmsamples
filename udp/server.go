package main

import (
	"net"
	"os"
	"os/signal"
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
	UDPAddr,err:=net.ResolveUDPAddr("udp",":8001")
	if err!=nil{
		panic(err)
	}
	conn,err:=net.ListenUDP("udp",UDPAddr)
	if err!=nil{
		panic(err)
	}
	stopCh:=make(chan os.Signal,1)
	signal.Notify(stopCh, os.Interrupt)
	go listenLoop(conn)
	<-stopCh
}