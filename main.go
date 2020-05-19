package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pub", pubKeyHandler)
	mux.HandleFunc("/sign", saveKeyHandler)
	lis, err := net.Listen("tcp", ":23333")
	if err != nil {
		fmt.Println("failed to listen 23333", err)
		return
	}
	panic(http.Serve(lis, mux))
}

func pubKeyHandler(writer http.ResponseWriter, _ *http.Request) {
	home := os.Getenv("HOME")
	key := fmt.Sprintf("%s/.ssh/id_ed25519.pub", home)
	log.Printf("start to read %s", key)
	fd, err := os.Open(key)
	if err != nil {
		log.Printf("failed to open file %v\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("failed to read all data %v\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func saveKeyHandler(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	home := os.Getenv("HOME")
	key := fmt.Sprintf("%s/.ssh/id_ed25519-cert.pub", home)
	log.Printf("begin to write data into %s", key)
	fd, err := os.OpenFile(key, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fd.Close()
	_, err = fd.Write(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
}
