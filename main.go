package main

import "net/http"

func main() {
	s := http.NewServeMux()
	s.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	s.HandleFunc("/healthz", readinessEndpoint)
	serv := new(http.Server)
	serv.Handler = s
	serv.Addr = ":8080"
	serv.ListenAndServe()
}

func readinessEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("OK"))
}
