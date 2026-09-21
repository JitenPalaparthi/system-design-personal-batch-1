package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var mu sync.Mutex
var items = map[int]Item{1: {ID: 1, Name: "Go item"}}
var nextID = 2

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func itemID(path string) (int, error) { return strconv.Atoi(strings.TrimPrefix(path, "/items/")) }

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	switch r.Method {
	case http.MethodGet:
		out := make([]Item, 0, len(items))
		for _, v := range items {
			out = append(out, v)
		}
		writeJSON(w, 200, out)
	case http.MethodPost:
		var in Item
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Name == "" {
			writeJSON(w, 400, map[string]string{"error": "name is required"})
			return
		}
		in.ID = nextID
		nextID++
		items[in.ID] = in
		writeJSON(w, 201, in)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func itemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := itemID(r.URL.Path)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid id"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	item, ok := items[id]
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, 200, item)
	case http.MethodPut:
		var in Item
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Name == "" {
			writeJSON(w, 400, map[string]string{"error": "name is required"})
			return
		}
		item.Name = in.Name
		items[id] = item
		writeJSON(w, 200, item)
	case http.MethodDelete:
		delete(items, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func main() {
	instance := os.Getenv("INSTANCE_NAME")
	if instance == "" {
		instance = "go-service"
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok", "service": "go", "instance": instance})
	})
	http.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"service": "go", "instance": instance})
	})
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)
	fmt.Println("Go service listening on :8080 as", instance)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
