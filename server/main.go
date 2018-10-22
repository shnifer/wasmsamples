package main

import (
	"net/http"
	"log"
	"time"
	"strconv"
)

type sessionData struct{
	user string
}

var m map[int64] sessionData

func init(){
	m=make(map[int64] sessionData)
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

func rootHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler")
	http.Redirect(w,r,"/login", http.StatusPermanentRedirect)
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
		http.Redirect(w,r,"/hello", http.StatusOK)
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
	}
}

func helloHandler (w http.ResponseWriter, r *http.Request) {
	log.Println("helloHandler")

	sessData:=checkSessionCookie(r)
	empty:=sessionData{}
	if sessData==empty{
		log.Println("somehow hello called without cookie")
		http.Redirect(w,r,"/login", 403)
		return
	}

	user := sessData.user
	log.Println("hello user ",user)
	w.Write([]byte("Hello "+user+"!"))
}

func main(){
//	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", loginHandler)
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		panic(err)
	}
}
