package main

import (
	"net/http"
	"strconv"
	"github.com/shnifer/nigiri/vec2"
	"sync"
	"time"
	"encoding/json"
	"log"
)

var iterMu *sync.Mutex
var idnIter int

var dataMu *sync.Mutex
var ships ShipsData
var lastSeen map[int]time.Time

func nextIdn() int{
	iterMu.Lock()
	defer iterMu.Unlock()

	idnIter++
	return idnIter
}

func init(){
	iterMu = &sync.Mutex{}
	dataMu = &sync.Mutex{}
	ships = make(ShipsData,0)
	lastSeen = make(map[int]time.Time)
}


func pingHandler(w http.ResponseWriter, r *http.Request) {
	dataMu.Lock()
	defer dataMu.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	idn:=idn(r)
	if idn==-1{
		newidn:=nextIdn()
		idnstr:=strconv.Itoa(newidn)
		name:=r.Header.Get("username")
		if name==""{
			name="<noname>"
		}
		name = "["+idnstr+"] "+name
		ships = append(ships, ShipData{
			Name: name,
			Idn: newidn,
		})

		w.Header().Set("idn", idnstr)
		return
	}

	buf,err:=json.Marshal(ships)
	if err!=nil{
		log.Println("can't marshal ships")
		return
	}
	w.Write(buf)
}

func gameListner(){
	mux:=http.NewServeMux()
	mux.HandleFunc("/",pingHandler)
	srv:=&http.Server{
		Addr:":8100",
		Handler:mux,
	}
	err:=srv.ListenAndServe()
	if err!=nil{
		panic(err)
	}
}

func idn(r *http.Request) int{
	str:=r.Header.Get("idn")
	i,err:=strconv.Atoi(str)
	if err!=nil{
		return -1
	}
	return i
}


type ShipData struct{
	Idn int
	Position vec2.V2
	Angle float64
	Name string
}

type ShipsData []ShipData