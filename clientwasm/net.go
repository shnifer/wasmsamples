package main

import (
	"github.com/shnifer/nigiri/vec2"
	"sync"
	"time"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type ShipData struct{
	Idn int
	Position vec2.V2
	Angle float64
	Name string
}

type ShipsData []ShipData

var mu *sync.Mutex
var Ships ShipsData
var client *http.Client

var myIdn int

func init(){
	Ships = make(ShipsData,0)
	mu = &sync.Mutex{}
	client = &http.Client{
		Timeout:time.Second/5,
	}
}

func networkLoop(){
	for{
		doRequest()
		time.Sleep(time.Second/10)
	}
}

func doRequest(){
	req,err:=http.NewRequest(http.MethodGet, gameHostName, nil)
	if err!=nil{
		log.Println("doRequest ",err)
		return
	}
	if myIdn==0{
		req.Header.Set("username", userName)
		resp, err:=client.Do(req)
		if err!=nil{
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		idnstr:=resp.Header.Get("idn")
		idn,err:=strconv.Atoi(idnstr)
		if err!=nil{
			log.Println(err)
		}
		myIdn = idn
		return
	}
	req.Header.Set("idn", strconv.Itoa(myIdn))
	resp, err:=client.Do(req)
	if err!=nil{
		log.Println("req do err: ",err)
		return
	}
	defer resp.Body.Close()

	buf, err:= ioutil.ReadAll(resp.Body)
	if err!=nil{
		log.Println("can't read buf ",err)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	err=json.Unmarshal(buf, &Ships)
	if err!=nil{
		log.Println("can't unmarshal ",err)
		return
	}
}