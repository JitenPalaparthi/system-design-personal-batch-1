package main

import (
  "encoding/json"
  "net/http"
  "os"
)

type Response struct { Service string `json:"service"`; Hostname string `json:"hostname"`; Path string `json:"path"` }
func reply(w http.ResponseWriter, r *http.Request) { h,_:=os.Hostname(); w.Header().Set("Content-Type","application/json"); json.NewEncoder(w).Encode(Response{"go-service",h,r.URL.Path}) }
func main(){
  http.HandleFunc("/", reply)
  http.HandleFunc("/go", reply)
  http.HandleFunc("/health", func(w http.ResponseWriter,r *http.Request){w.Write([]byte("OK"))})
  http.ListenAndServe(":8080",nil)
}
