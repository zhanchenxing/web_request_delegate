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
}

func main(){
	fmt.Println("Hello, world!")
	conn, err := net.Dial("tcp", "localhost:8080" )
	checkError(err)

	reader := bufio.NewReader( conn )
	start_time := time.Now()

	results := make( chan string, 12 )
	count := 0

	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF{
			fmt.Println("all readed!")

			for i := 0; i < count; i++{
				result := <-results
				fmt.Println("result:", result)
			}

			end_time := time.Now()

			last_time := end_time.Sub(start_time)
			fmt.Println("time used:", last_time)

			os.Exit(0)
		}

		checkError(err)

		fmt.Println("line:", string(line) )
		fmt.Println("isPrefix:", isPrefix)

		count ++
		//_, err = http.Get(string(line))
		//fmt.Println("error", err, ", url:", string(line))
		go fetch( string(line), results)

		conn.Write([]byte("ok"))
	}
}

func checkError( err error ){
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
