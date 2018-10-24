package main

import (
	"github.com/shnifer/nigiri"
	"github.com/hajimehoshi/ebiten"
	"net/http"
	"image"
	_"image/png"
	"golang.org/x/image/font"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"log"
	"github.com/shnifer/nigiri/vec2"
)

var shipSprite *nigiri.Sprite
var shipText *nigiri.TextSprite
var textFace font.Face

func loader(name string) (*ebiten.Image, error){
	log.Println(webHostName+"/"+name)
	resp, err:=http.Get(webHostName+"/"+name)
	if err!=nil{
		return nil, err
	}
	defer resp.Body.Close()

	img,_,err:=image.Decode(resp.Body)
	if err!=nil{
		return nil,err
	}
	ebiImg,_:=ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return ebiImg, nil
}
func faceLoader (name string, size float64) (font.Face, error) {
	resp, err:=http.Get(webHostName+"/"+name)
	if err!=nil{
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	face, err := newFace(dat, size)
	if err != nil {
		return nil, err
	}

	return face, nil
}


func initGraph(){
	nigiri.SetTexLoader(loader)
	shipTex,err := nigiri.GetTex("ship.png")
	if err!=nil{
		panic(err)
	}
	nigiri.SetFaceLoader(faceLoader)
	face,err:=nigiri.GetFace("furore.ttf", 10)
	if err!=nil{
		panic(err)
	}
	shipSprite = nigiri.NewSprite(shipTex, 0)
	shipSprite.Pivot = vec2.Center

	shipText = nigiri.NewTextSprite(1, false, 1)
	shipText.ChangeableTex = true
	shipText.DefFace = face
}

func newFace(b []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	tto := &truetype.Options{
		Size: size,
	}
	face := truetype.NewFace(f, tto)
	return face, nil
}