package main

import (
	"github.com/shnifer/nigiri"
	"github.com/hajimehoshi/ebiten"
	"log"
	"net/url"
	"github.com/gopherjs/gopherwasm/js"
)

var webHostName string
var gameHostName string
var userName string
var Q *nigiri.Queue

func mainLoop(win *ebiten.Image, dt float64) error{
	if ebiten.IsDrawingSkipped(){
		return  nil
	}
	mu.Lock()
	defer mu.Unlock()

	Q.Clear()
	for _,dat:=range Ships{
		shipSprite.Position = dat.Position
		shipSprite.Angle = dat.Angle
		Q.Add(shipSprite)

		shipText.Position = dat.Position
		shipText.SetText(dat.Name)
		Q.Add(shipText)
	}
	Q.Run(win)
	return nil
}

func main(){
	Q = nigiri.NewQueue()
	w,h := 500,500
	if nigiri.IsJS {
		w, h = ebiten.ScreenSizeInFullscreen()
		ebiten.SetFullscreen(true)
	}
	getHostName()
	initGraph()
	go networkLoop()
	nigiri.Run(mainLoop,w,h,1,"test")
}

func getHostName(){
	if nigiri.IsJS{
		urlstr:=js.Global().Get("document").Get("URL").String()
		log.Println("url: ", urlstr)
		ur,err:=url.Parse(urlstr)
		if err!=nil{
			log.Println("error parsing url: ", err)
			return
		}
		webHostName="http://"+ur.Hostname()+":8080"
		gameHostName="http://"+ur.Hostname()+":8100"
		userName="get me from cookies"

	} else {
		webHostName="http://localhost:8080"
		gameHostName="http://localhost:8100"
		userName="testLocal"
	}
}