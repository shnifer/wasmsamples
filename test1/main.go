package main

import (
	"log"
	"github.com/hajimehoshi/ebiten"
)

func loop(w *ebiten.Image) error{
	log.Println("Hello!")
	return nil
}

func main(){
	ebiten.Run(loop, 100,100,1,"test")
}
