package main

import (
	"fmt"
)

func panic_conn(conn_obj *tcp_conn_object, reason string) {

	if conn_obj.conn == nil {
		return
	}

	conn_obj.conn.Close()

	log_err("PANIC CLOSE CONN_ID: " + conn_obj.conn_id + " TYPE: " + fmt.Sprint(conn_obj.conn_type) + " REASON: " + reason)

	remove_conn(conn_obj)

	//this doesn't exactly mean that SOCKET is down (BUT CONN IS DOWN)
	//need contact with other failover routers to confirm
	//create logic to decide how to reestablish connection and update routing tables
}

func socket_error(conn_obj *tcp_conn_object) {

	// REROUTE SOCKET'S CACHE

}
