package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var count int
var born time.Time

func handler(w http.ResponseWriter, r *http.Request) {

	myname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hello from container %q\n", myname)
	count++
	fmt.Fprintf(w, "I have been visted %d times\n", count)
	fmt.Fprintf(w, "I was born %q \n", born.Format(time.UnixDate))
}

func main() {
	born = time.Now()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
