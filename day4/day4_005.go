// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// get time now
	start := time.Now()
	// !! keyword chan is required for parallel programming of Golang
	// make(chan string) for transfer string between different thread
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		// parallel keyword: go
		go fetch(url, ch) // start a goroutine
	}
	// 10 urls => 10 print
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	// time.Since(start).Seconds() get costs of time
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
func fetch(url string, ch chan<- string) {
	start := time.Now()

	// set timeout 5s
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		// Sprint func(a ...interface{}) string
		// Sprint formats using the default formats for its operands and returns the resulting string. Spaces are added between operands when neither is a string.
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	// ioutil.Discard is an io.Writer on which all Write calls succeed without doing anything.
	// io.Copy returns (written int64, err error)
	nbytes, err := io.Copy(os.Stdout, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	// get time costs of call 1 time fetch function
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
