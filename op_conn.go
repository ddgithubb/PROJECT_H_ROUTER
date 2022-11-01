package main

import (
	"math"
	"net"
	"strings"
	"sync/atomic"
	"unsafe"
)

func add_new_conn_route(conn_id string, conn_group string, add_conn *net.TCPConn) (conn_obj *tcp_conn_object, ok bool) {

	conn_obj = &tcp_conn_object{
		conn:       add_conn,
		conn_id:    conn_id,
		conn_group: conn_group,
	}
	ok = false

	if strings.Contains(conn_id, "SOCK_") {

		var cache cache_writer

		if add_conn != nil {
			// CREATE or OPEN FILE
		}
		conn_obj.conn_type = SOCKET_TYPE

		socket_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&socket_table.p))
		_, exist := table_old[conn_id]

		if !exist {
			table := make(conn_table)
			for k, v := range table_old {
				table[k] = v
			}
			table[conn_id] = &conn_route{
				conn_obj: conn_obj,
				cache:    &cache,
			}
			atomic.StorePointer(&socket_table.p, unsafe.Pointer(&table))
			ok = true
		}

		socket_table.Unlock()

	} else if strings.Contains(conn_id, "API_") {

		conn_obj.conn_type = API_TYPE

		api_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&api_table.p))
		_, exist := table_old[conn_id]

		if !exist {
			table := make(conn_table)
			for k, v := range table_old {
				table[k] = v
			}
			table[conn_id] = &conn_route{
				conn_obj: conn_obj,
				cache:    nil,
			}
			atomic.StorePointer(&api_table.p, unsafe.Pointer(&table))
			ok = true
		}

		api_table.Unlock()

	} else if strings.Contains(conn_id, "ROUTER_") {
		conn_obj.conn_type = ROUTER_TYPE
	}

	return
}

func recover_conn_route(conn_id string, add_conn *net.TCPConn) (conn_obj *tcp_conn_object, ok bool) {

	conn_obj = &tcp_conn_object{
		conn:       add_conn,
		conn_id:    conn_id,
		conn_type:  math.MaxUint8,
		conn_group: "",
	}
	ok = false

	if strings.Contains(conn_id, "SOCK_") {

		socket_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&socket_table.p))
		existing_route, exist := table_old[conn_id]

		if exist {
			table := make(conn_table)
			for k, v := range table_old {
				table[k] = v
			}
			table[conn_id].conn_obj.conn = add_conn
			atomic.StorePointer(&socket_table.p, unsafe.Pointer(&table))
		}

		socket_table.Unlock()

		if exist {
			conn_obj = existing_route.conn_obj
			conn_obj.conn = add_conn
			ok = true
		}

	} else if strings.Contains(conn_id, "API_") {

		api_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&api_table.p))
		existing_route, exist := table_old[conn_id]

		if exist {
			table := make(conn_table)
			for k, v := range table_old {
				table[k] = v
			}
			table[conn_id].conn_obj.conn = add_conn
			atomic.StorePointer(&api_table.p, unsafe.Pointer(&table))
		}

		api_table.Unlock()

		if exist {
			conn_obj = existing_route.conn_obj
			conn_obj.conn = add_conn
			ok = true
		}

	} else if strings.Contains(conn_id, "ROUTER_") {
		conn_obj.conn_type = ROUTER_TYPE
	}

	return
}

func remove_conn(conn_obj *tcp_conn_object) {

	switch conn_obj.conn_type {
	case SOCKET_TYPE:
		socket_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&socket_table.p))

		table := make(conn_table)
		for k, v := range table_old {
			if v.conn_obj.conn_id != conn_obj.conn_id {
				table[k] = v
			}
		}

		atomic.StorePointer(&socket_table.p, unsafe.Pointer(&table))

		socket_table.Unlock()
	case API_TYPE:
		api_table.Lock()

		table_old := *(conn_table_ptr)(atomic.LoadPointer(&api_table.p))

		table := make(conn_table)
		for k, v := range table_old {
			if v.conn_obj.conn_id != conn_obj.conn_id {
				table[k] = v
			}
		}

		atomic.StorePointer(&api_table.p, unsafe.Pointer(&table))

		api_table.Unlock()
	case ROUTER_TYPE:
	}

}

// func purge_conn_route(conn_id string) {

// }
