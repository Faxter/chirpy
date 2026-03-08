package main

import "net/http"

func main() {
	s := http.NewServeMux()
	s.Handle("/", http.FileServer(http.Dir(".")))
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}
