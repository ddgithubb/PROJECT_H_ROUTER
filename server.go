package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"time"
)

//CONFIG (to be put in config file)
const (
	VALID_NANOID_CHAR string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	DEFAULT_TIMEOUT_AFTER_WRITE time.Duration = time.Second * 1
	MAX_RETRIES                               = 10

	CLIENT_REQUEST_TIMEOUT time.Duration = time.Second * 3

	HEARTBEAT_MAX_DELAY time.Duration = time.Second * 10 // This includes network/buffer congestion

	CONCURRENCY_FACTOR = 1

	BUFFER_SIZE = 1000000 //1 mb
)

const SOCKET_GROUP = "AZ1_MAIN"

var router_addr = "127.0.0.1:10000"

var CONCURRENCY uint32 = uint32(runtime.NumCPU() * CONCURRENCY_FACTOR)

var router_logger *log.Logger

//var op_pool *ants.PoolWithFunc

func init() {
	router_logs_file, err := os.OpenFile("router_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	router_logger = log.New(router_logs_file, "", log.LstdFlags)
}

func main() {

	// op_pool, err := ants.NewPoolWithFunc(ants.DefaultAntsPoolSize, handle_op)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// defer op_pool.Release()

	go func() {
		_ = pprof.Handler("")
		http.ListenAndServe("localhost:10001", nil)
	}()

	laddr, err := net.ResolveTCPAddr("tcp", router_addr)
	if err != nil {
		log.Fatalln(err)
	}

	router, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Listening on addr: ", router_addr)

	for {
		conn, err := router.AcceptTCP()
		if err != nil {
			log.Println(err)
		}

		go new_conn(conn)
	}

}

func log_err(err string) {
	router_logger.Println(err)
}
