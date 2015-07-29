package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
	"log"
	"math/rand"
	"net/http"
)

const PORT = ":3000"

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "redis.default.cluster.local:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println("[CONNECTED] Redis")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleCreate).Methods("POST")
	r.HandleFunc("/{id}", handleFind).Methods("GET")
	log.Println("[STARTED] on port " + PORT)
	http.ListenAndServe(PORT, r)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	if len(url) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "missing url param")
		return
	}
	exists := true
	var code string
	var err error
	for exists {
		code = string(randCode(7))

		exists, err = client.Exists(code).Result()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, err.Error())
			return
		}
	}
	err = client.Set(code, url, 0).Err()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, code)
}

func handleFind(w http.ResponseWriter, r *http.Request) {
	url, err := client.Get(mux.Vars(r)["id"]).Result()
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "not found")
		return
	}
	http.Redirect(w, r, url, 301)
}

func randCode(length int) []byte {
	src := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = src[rand.Intn(len(src))]
	}
	return result
}
