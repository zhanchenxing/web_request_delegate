package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"io"
	"time"
	"net/http"
)

func fetch( url string, to chan string ) {
	_, err := http.Get(url)
	fmt.Println("error:", err, ", url:", url)
	to <- "ok"
}

func read_url( from net.Conn, to chan string ){
	reader := bufio.NewReader( from )
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF{
			fmt.Println("all readed!")
			close(to)
			return
		}
		fmt.Println(string(line), isPrefix, err)
		checkError(err)
		to <-string(line)
	}
}

func main(){
	fmt.Println("Hello, world!")
	conn, err := net.Dial("tcp", "localhost:8080" )
	checkError(err)

	start_time := time.Now()

	results := make( chan string, 12 )
	url_chan := make( chan string, 12 )
	go read_url( conn, url_chan )

	count := 0

	loop:
	for {
		select {
		case url, ok := <-url_chan:
			if ok{
				count ++
				go fetch( url, results)
			} else {
				fmt.Println("url_chan closed!")
				break loop
			}
		case result := <-results:
			count --
			fmt.Println("result:", result)
			conn.Write([]byte(result))
		}
	}

	fmt.Println("remaining url count:", count)
	for i:=0; i<count; i++{
		result := <-results
		fmt.Println("result:", result)
		conn.Write([]byte(result))
	}


	end_time := time.Now()
	last_time := end_time.Sub(start_time)
	fmt.Println("time used:", last_time)
}

func checkError( err error ){
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
