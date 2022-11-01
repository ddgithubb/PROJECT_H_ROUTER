package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

func read_tcp(conn_obj *tcp_conn_object) {

	var received uint32

	heartbeat_chan := make(chan []byte, 3)

	defer func() {
		close(heartbeat_chan)
	}()

	switch conn_obj.conn_type {
	case SOCKET_TYPE:
		go sender_heartbeat(conn_obj, heartbeat_chan)
	case API_TYPE:
		go receiver_heartbeat(conn_obj, &received, heartbeat_chan)
	case ROUTER_TYPE:
	default:
	}

	reader := bufio.NewReaderSize(conn_obj.conn, BUFFER_SIZE)

	var (
		err    error
		header []byte = make([]byte, 5)
		size   uint32
		op     byte
		b      []byte
	)
	for {

		_, err = io.ReadFull(reader, header)
		if err != nil {
			log_err(err.Error())
			break
		}

		op = header[0]

		size = binary.BigEndian.Uint32(header[1:])
		b = nil

		if size > 0 {
			b = make([]byte, size)
			_, err = io.ReadFull(reader, b)
			if err != nil {
				log_err(err.Error())
				break
			}
		}

		fmt.Println("RECV:", header, b[0], string(b[1:]))

		if op == 1 {
			heartbeat_chan <- b

			if conn_obj.conn_type == API_TYPE {
				atomic.StoreUint32(&received, 0)
			}
		} else {
			switch conn_obj.conn_type {
			case SOCKET_TYPE:
				handle_socket_op(conn_obj, b, op)
			case API_TYPE:
				atomic.AddUint32(&received, 1)
				handle_api_op(conn_obj, b, op)
			case ROUTER_TYPE:
			}
		}

	}

}

func sender_heartbeat(conn_obj *tcp_conn_object, heartbeat_chan chan []byte) {

	var b []byte
	var ok bool
	var cur_ver byte = 0
	var timeout_ver byte
	var received uint32
	timeout_chan := make(chan byte)

	for {
		select {
		case b, ok = <-heartbeat_chan:

			if !ok {
				return
			}

			if byte(cur_ver+1) != b[0] {
				panic_conn(conn_obj, "heartbeat mismatch "+fmt.Sprint(byte(cur_ver+1))+" vs "+fmt.Sprint(b[0]))
				return
			}

			cur_ver = b[0]

			//fmt.Println("Heartbeat from", conn_obj.conn_id, cur_ver)

			write(conn_obj, []byte{1, 0, 0, 0, 1, cur_ver})

			received = binary.BigEndian.Uint32(b[9:13])

			if received != 0 {
				fmt.Println("Amount sent to socket:", received)
			}

			// route, ok := (*(conn_table_ptr)(atomic.LoadPointer(&socket_table.p)))[conn_obj.conn_id]
			// if !ok {
			// 	log_err("sock_id: " + conn_obj.conn_id + " requested during heartbeat but not found")
			// 	return
			// }

			// route.cache.Lock()
			// DO STUFF WITH received
			// route.cache.Unlock()

			go func(prev_ver byte, b []byte) {
				next_expected_heartbeat_unix_nano := int64(binary.BigEndian.Uint64(b[1:9]))
				time.Sleep(time.Duration(next_expected_heartbeat_unix_nano-time.Now().UnixNano()) + HEARTBEAT_MAX_DELAY)
				timeout_chan <- prev_ver
			}(cur_ver, b)
		case timeout_ver = <-timeout_chan:
			if cur_ver == timeout_ver {
				panic_conn(conn_obj, "timed out ver "+fmt.Sprint(timeout_ver)+" vs "+fmt.Sprint(cur_ver))
				return
			}
		}
	}

}

func receiver_heartbeat(conn_obj *tcp_conn_object, received *uint32, heartbeat_chan chan []byte) {

	var b []byte
	var ok bool
	var cur_ver byte = 0
	var timeout_ver byte
	timeout_chan := make(chan byte)
	heartbeat_packaged_op := make([]byte, 10)

	_ = copy(heartbeat_packaged_op[:5], []byte{1, 0, 0, 0, 5})

	for {
		select {
		case b, ok = <-heartbeat_chan:

			if !ok {
				return
			}

			if byte(cur_ver+1) != b[0] {
				panic_conn(conn_obj, "heartbeat mismatch "+fmt.Sprint(byte(cur_ver+1))+" vs "+fmt.Sprint(b[0]))
				return
			}

			cur_ver = b[0]

			//fmt.Println("Heartbeat from", conn_obj.conn_id, cur_ver)

			heartbeat_packaged_op[5] = cur_ver
			binary.BigEndian.PutUint32(heartbeat_packaged_op[6:], atomic.LoadUint32(received))
			write(conn_obj, heartbeat_packaged_op)

			go func(prev_ver byte, b []byte) {
				next_expected_heartbeat_unix_nano := int64(binary.BigEndian.Uint64(b[1:9]))
				time.Sleep(time.Duration(next_expected_heartbeat_unix_nano-time.Now().UnixNano()) + HEARTBEAT_MAX_DELAY)
				timeout_chan <- prev_ver
			}(cur_ver, b)
		case timeout_ver = <-timeout_chan:
			if cur_ver == timeout_ver {
				panic_conn(conn_obj, "timed out ver "+fmt.Sprint(timeout_ver)+" vs "+fmt.Sprint(cur_ver))
				return
			}
		}
	}

}
