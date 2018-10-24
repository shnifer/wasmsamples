package main

import (
	"net/http"
	"strconv"
	"github.com/shnifer/nigiri/vec2"
	"sync"
	"time"
	"encoding/json"
	"log"
	"io/ioutil"
	"math"
)

var iterMu *sync.Mutex
var idnIter int

var dataMu *sync.Mutex
var ships ShipsData
var lastSeen map[int]time.Time
var targets map[int]vec2.V2

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
	targets = make(map[int]vec2.V2)
}

func sendOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method==http.MethodOptions{
		sendOptions(w,r)
		return
	}
	dataMu.Lock()
	defer dataMu.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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

		w.Write([]byte(idnstr))
		log.Println("send to user ", name, "id=",idnstr)
		return
	}

	buf,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		log.Println(err)
	}
	var target vec2.V2
	err=json.Unmarshal(buf, &target)
	if err!=nil{
		log.Println(err)
	}
	targets[idn] = target
	lastSeen[idn] = time.Now()

	defer r.Body.Close()

	buf,err=json.Marshal(ships)
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

func GameCycle(){
	for{
		gameTick()
		time.Sleep(time.Second/10)
	}
}

const dt=0.1
const rotSpeed = 120
const moveSpeed = 100

func gameTick(){
	dataMu.Lock()
	defer dataMu.Unlock()

	var del []int

	now:=time.Now()
	for id, t:=range lastSeen{
		if now.Sub(t).Seconds()>1{
			del = append(del, id)
		}
	}
	for id:=range del{
		delete(lastSeen, id)
		delete(targets, id)
		for j:=range ships{
			if ships[j].Idn==id{
				ships = append(ships[:j], ships[j+1:]...)
				break
			}
		}
	}

	for i, ship:=range ships{
		id:=ship.Idn
		vec:=targets[id].Sub(ship.Position)
		dir:=vec.Dir()
		l:=vec.Len()
		if l==0{
			continue
		}
		dAng:=dir-ship.Angle
		if math.Abs(dAng)>1 {
			maxChange:=dt*rotSpeed
			if dAng>maxChange {
				dAng = maxChange
			} else if dAng< (-maxChange){
				dAng = -maxChange
			}
			ship.Angle +=dAng
			ships[i]=ship
			continue
		}
		maxMove:=dt*moveSpeed
		if l>maxMove{
			l=maxMove
		}
		ship.Position = ship.Position.AddMul(vec.Normed(),l)
		ships[i] = ship
	}
}