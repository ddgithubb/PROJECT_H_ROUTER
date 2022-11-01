package main

func get_prev_state(prev_sock_id string, prev_ws_id string, new_sock_id string, new_ws_id string) {

	if sent := route_to_ws(prev_sock_id, prev_ws_id, 30, []string{new_sock_id, new_ws_id}, nil); !sent {
		send_prev_state(new_sock_id, new_ws_id, "0", []byte("{}"))
	}

}

func send_prev_state(new_sock_id string, new_ws_id string, status string, state []byte) {

	route_to_ws(new_sock_id, new_ws_id, 31, []string{status}, state)

}
