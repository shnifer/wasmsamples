package main

import (
	"log"
	"github.com/gopherjs/gopherwasm/js"
	"net/http"
)

func main(){
	addr:=js.Global().Get("document").Get("url")
	log.Println("we are at the ",addr)
	http.Request{}.ParseForm()
}
