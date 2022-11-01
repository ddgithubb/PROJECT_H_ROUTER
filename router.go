package main

import (
	"errors"
	"net"
	"sync/atomic"
)

// Writes exclusively to specific conn
// HELPER FUNCTION FOR ONLY router.go
func write(conn_obj *tcp_conn_object, packaged packaged_op) error {

	var err error

	if conn_obj.conn == nil {
		return errors.New("unable to write to nil conn")
	}

	for i := 0; i < MAX_RETRIES; i++ {
		_, err = conn_obj.conn.Write(packaged)
		if err == nil {
			break
		}
		if !err.(*net.OpError).Temporary() {
			break
		}
	}

	if err != nil {
		panic_conn(conn_obj, "write faliure: "+err.Error())
	}

	return err
}

func write_op(conn_obj *tcp_conn_object, op byte, params []string, payload []byte) error {
	return write(conn_obj, package_op(op, params, payload))
}

// Routes directly to specified socket, then to the specific ws
// CRITICAL to succeed (tries anything to get to ws)
func route_to_ws_direct(conn_obj *tcp_conn_object, ws_id string, underlying_op byte, underlying_params []string, underlying_payload []byte) error {
	return write(conn_obj, package_op_router_wrapper(50, []string{ws_id}, underlying_op, underlying_params, underlying_payload))
}

// Routes payload to specified socket, then to the specific ws
// TRIES to succeed (tries anything to get to destination)
func route_to_ws(sock_id string, ws_id string, underlying_op byte, underlying_params []string, underlying_payload []byte) bool {
	route, ok := (*(conn_table_ptr)(atomic.LoadPointer(&socket_table.p)))[sock_id]
	if !ok {
		return false
	}

	if route.conn_obj.conn != nil {
		err := write(route.conn_obj, package_op_router_wrapper(50, []string{ws_id}, underlying_op, underlying_params, underlying_payload))
		if err != nil {
			socket_error(route.conn_obj) // already panic conned
		}
	} else {
		// TODO: Rerouting strategy (not cached based)
	}

	return true
}

// Routes payload to specified socket, then to the specific ws. Cached
// CRITICAL to succeed (tries anything to get to destination)
func route_message_to_ws(sock_id string, ws_id string, underlying_op byte, underlying_params []string, underlying_payload []byte) bool {
	route, ok := (*(conn_table_ptr)(atomic.LoadPointer(&socket_table.p)))[sock_id]
	if !ok {
		return false
	}

	if route.conn_obj.conn != nil {

		// route.cache.Lock()

		// route.cache.writer.Write(packaged_op)
		err := write(route.conn_obj, package_op_router_wrapper(51, []string{ws_id}, underlying_op, underlying_params, underlying_payload))

		// route.cache.Unlock()

		if err != nil {
			socket_error(route.conn_obj) // already panic conned
		}

	} else {
		// TODO: Rerouting strategy
	}

	return true
}

////////////////////////////////////////////////////////////////
////////////////////////// PLANNED /////////////////////////////
////////////////////////////////////////////////////////////////

// // Routes payload to specified socket (is not conn specific)
// // CRITICAL to succeed (tries anything to get to socket)
// func route_to_sock_out(sock_id string, underlying_op byte, underlying_params []string, underlying_payload []byte) {
// 	route, ok := (*(conn_table_ptr)(atomic.LoadPointer(&socket_table.p)))[sock_id]
// 	if !ok {
// 		log_err("sock_id: " + sock_id + " requested but not found")
// 		return
// 	}

// 	packaged := package_op(underlying_op, underlying_params, underlying_payload)

// 	if route.conn_obj.conn != nil {
// 		err := write(route.conn_obj, packaged)
// 		if err != nil {
// 			// TODO: Rerouting strategy
// 		}
// 	} else {
// 		// TODO: Rerouting strategy
// 	}
// }

// func route_to_sock_in(sock_id string, packaged []byte) {
// 	route, ok := (*(conn_table_ptr)(atomic.LoadPointer(&socket_table.p)))[sock_id]
// 	if !ok {
// 		log_err("sock_id: " + sock_id + " requested but not found")
// 		return
// 	}

// 	if route.conn_obj.conn != nil {
// 		err := write(route.conn_obj, packaged)
// 		if err != nil {
// 			// TODO: Rerouting strategy
// 		}
// 	} else {
// 		// TODO: Rerouting strategy
// 	}
// }

// func route_from_router_to_ws_out(sock_id string, underlying_op byte, underlying_params []string, underlying_payload []byte) {

// }

// func route_from_router_to_ws_in(sock_id string, ws_id string, underlying_payload []byte) {

// }
