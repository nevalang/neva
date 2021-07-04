package api

import (
	"log"
	"net/http"
)

func ListenAndServe(addr string) error {
	log.Println("Api server started on " + addr)
	return http.ListenAndServe(addr, nil)
}
