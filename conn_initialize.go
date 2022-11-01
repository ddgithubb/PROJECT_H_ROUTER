package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func new_conn(conn *net.TCPConn) {

	defer conn.Close()

	err := conn.SetReadDeadline(time.Now().Add(DEFAULT_TIMEOUT_AFTER_WRITE))
	if err != nil {
		log_err(err.Error())
		return
	}

	header := make([]byte, 5)
	_, err = conn.Read(header)
	if err != nil {
		log_err(err.Error())
		return
	}

	size := binary.BigEndian.Uint32(header[1:5])

	b := make([]byte, size)
	_, err = conn.Read(b)
	if err != nil {
		log_err(err.Error())
		return
	}

	var conn_obj *tcp_conn_object
	var ok bool

	switch header[0] {
	case 2:

		params, _, err := byte_to_params_and_payload(b, 3, false)
		if err != nil {
			log_err(err.Error())
			return
		}

		if params[0] == "" || params[1] == "" {
			log_err("Expected conn_id/conn_group: but got ''")
			return
		}

		// if params[2] != "nil" {
		// 	purge_conn_route(params[2])
		// }

		conn_obj, ok = add_new_conn_route(params[0], params[1], conn)
		if err != nil {
			return
		}

		if ok {
			err = write_op(conn_obj, 2, []string{"1"}, nil)
			if err != nil {
				return
			}
		} else {
			err = write_op(conn_obj, 2, []string{"0"}, nil)
			return
		}

	case 3:

		conn_id := string(b)
		if conn_id == "" {
			log_err("Expected conn_id: but got ''")
			return
		}
		fmt.Println(conn_id)

		conn_obj, ok = recover_conn_route(conn_id, conn)

		if ok {
			err = write_op(conn_obj, 3, []string{"1"}, nil)
			if err != nil {
				return
			}
		} else {
			err = write_op(conn_obj, 3, []string{"0"}, nil)
			return
		}

	default:
		log_err("Expected op:2,3. Instead got " + fmt.Sprint(b[0]))
		return
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log_err(err.Error())
		return
	}

	//add_conn_route to neighbouring router groups
	//Yes it is redudant as all routers within 1 router group will be applying this operation
	//But extra redundancy is good

	read_tcp(conn_obj)
}
