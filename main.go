package main

import "net/http"

func main() {
	s := new(http.ServeMux)
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}
