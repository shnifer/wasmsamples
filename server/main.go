package main

import (
	"net/http"
	"log"
	"strconv"
	"time"
	"net"
)

type sessionData struct{
	user string
}

var m map[int64] sessionData

func init(){
	m=make(map[int64] sessionData)
}

func rootHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler ",r.URL.String(), r.Method)
	sessData:=checkSessionCookie(r)
	if (sessData==sessionData{}){
		http.Redirect(w,r,"/login", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w,r,"/hello", http.StatusTemporaryRedirect)
}

func loginHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("loginHandler")
	if r.Method == http.MethodGet {
		sessData:=checkSessionCookie(r)
		empty:=sessionData{}
		if sessData==empty {
			http.ServeFile(w, r, "server/login.html")
			return
		}
		log.Println("coming user with name ",sessData.user)
		http.Redirect(w,r,"/hello", http.StatusTemporaryRedirect)
	}
	if r.Method == http.MethodPost{
		r.ParseForm()
		user:=r.Form.Get("username")
		pass:=r.Form.Get("password")
		valid, sessID:=checkLogin(user,pass)
		if !valid{
			http.Error(w, "not valid pair", http.StatusUnauthorized)
			return
		}
		cookieStr:=strconv.Itoa(int(sessID))
		log.Println("user-pass ",user, "-",pass," logged, got cookie ",cookieStr)
		cookie:=&http.Cookie{
			Name: "sessionID",
			Value: cookieStr,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w,r,"/hello", http.StatusTemporaryRedirect)
	}
}

func helloHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("helloHandler")

	sessData:=checkSessionCookie(r)
	empty:=sessionData{}
	if sessData==empty{
		log.Println("somehow hello called without cookie")
		http.Redirect(w,r,"/login", http.StatusTemporaryRedirect)
		return
	}

	http.ServeFile(w,r,"server/wasmwrap.html")
}

func fileHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("filehadler: "+"server"+r.URL.String())
	http.ServeFile(w,r,"server"+r.URL.String())
}

func main(){
	go gameListner()
	go GameCycle()
	go udpServer()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/wasm_exec.js", fileHandler)
	http.HandleFunc("/main.wasm", fileHandler)
	http.HandleFunc("/ship.png", fileHandler)
	http.HandleFunc("/furore.ttf", fileHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


func checkSessionCookie(r *http.Request) sessionData{
	log.Println("checkSession")
	cookie,err:=r.Cookie("sessionID")
	if err!=nil{
		return sessionData{}
	}
	uid,err:=strconv.Atoi(cookie.Value)
	if err!=nil{
		return sessionData{}
	}
	log.Println("uid cookie: ",uid)
	sd,ok:=m[int64(uid)]
	if !ok{
		return sessionData{}
	}
	return sd
}

func checkLogin(user, pass string) (isValid bool, uid int64) {
	if user=="" || pass == ""{
		return false, 0
	}
	//check nice pass, but this time any is ok
	uid = time.Now().UnixNano()
	m[uid] = sessionData{
		user: user,
	}
	return true, uid
}

const msgSize = 10000

func listenLoop(conn *net.UDPConn){
	buf := make([]byte, msgSize)
	for{
		n, addr, err:= conn.ReadFromUDP(buf)
		if err!=nil {
			log.Println("ReadFromUDP err: ",err)
			continue
		}
		log.Println("recieved",n,"bytes from", addr,": ",string(buf[:n]))
		resp:=append([]byte("pong: "), buf[:n]...)
		conn.WriteTo(resp, addr)
	}
}

func udpServer(){
	UDPAddr,err:=net.ResolveUDPAddr("udp",":8001")
	if err!=nil{
		panic(err)
	}
	conn,err:=net.ListenUDP("udp",UDPAddr)
	if err!=nil{
		panic(err)
	}
	listenLoop(conn)
}