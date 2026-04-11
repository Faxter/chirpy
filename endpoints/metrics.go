package endpoints

import "net/http"

func (a *ApiConfig) IncrementsMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		a.FileServerHits.Add(1)
		handler.ServeHTTP(writer, req)
	})
}
