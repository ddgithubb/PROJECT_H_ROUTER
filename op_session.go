package main

func connect_session(conn_obj *tcp_conn_object, user_id string, session_id string, ws_id string) {

	var prev_session *session_routing_info
	shard := routing_table.get_shard(user_id)

	shard.Lock()

	routes, ok := shard.table[user_id]
	if !ok {
		shard.table[user_id] = make([]*session_routing_info, 1)
		shard.table[user_id][0] = &session_routing_info{
			session_id: session_id,
			ws_id:      ws_id,
			sock_id:    conn_obj.conn_id,
			active:     true,
		}
	} else {
		ok = false
		for i := 0; i < len(routes); i++ {
			if routes[i].session_id == session_id {
				ok = true
				prev_session = &session_routing_info{
					ws_id:   routes[i].ws_id,
					sock_id: routes[i].sock_id,
				}
				shard.table[user_id][i].ws_id = ws_id
				shard.table[user_id][i].sock_id = conn_obj.conn_id
				shard.table[user_id][i].active = true
			}
		}
		if !ok {
			shard.table[user_id] = append(routes, &session_routing_info{
				session_id: session_id,
				ws_id:      ws_id,
				sock_id:    conn_obj.conn_id,
				active:     true,
			})
		}
	}

	if prev_session != nil {
		route_to_ws(prev_session.sock_id, prev_session.ws_id, 32, nil, nil)
	}

	shard.Unlock()

	route_to_ws_direct(conn_obj, ws_id, 20, []string{"1"}, nil)
}

func recover_session(conn_obj *tcp_conn_object, user_id string, session_id string, ws_id string, prev_ws_id string) {

	var prev_session *session_routing_info

	shard := routing_table.get_shard(user_id)

	shard.Lock()

	routes, ok := shard.table[user_id]
	if !ok {
		shard.table[user_id] = make([]*session_routing_info, 1)
		shard.table[user_id][0] = &session_routing_info{
			session_id: session_id,
			ws_id:      ws_id,
			sock_id:    conn_obj.conn_id,
			active:     true,
		}
	} else {
		ok = false
		for i := 0; i < len(routes); i++ {
			if routes[i].session_id == session_id {
				ok = true
				prev_session = &session_routing_info{
					ws_id:   routes[i].ws_id,
					sock_id: routes[i].sock_id,
				}
				shard.table[user_id][i].ws_id = ws_id
				shard.table[user_id][i].sock_id = conn_obj.conn_id
				shard.table[user_id][i].active = true
			}
		}
		if !ok {
			shard.table[user_id] = append(routes, &session_routing_info{
				session_id: session_id,
				ws_id:      ws_id,
				sock_id:    conn_obj.conn_id,
				active:     true,
			})
		}
	}

	shard.Unlock()

	if ok {
		if prev_session.ws_id == prev_ws_id {
			route_to_ws_direct(conn_obj, ws_id, 21, []string{"1"}, nil)
			get_prev_state(prev_session.sock_id, prev_session.ws_id, conn_obj.conn_id, ws_id)
			return
		} else {
			route_to_ws(prev_session.sock_id, prev_session.ws_id, 32, nil, nil)
		}
	}

	route_to_ws_direct(conn_obj, ws_id, 21, []string{"0"}, nil)

}

func remove_session(conn_obj *tcp_conn_object, user_id string, session_id string, ws_id string) {

	shard := routing_table.get_shard(user_id)

	shard.Lock()

	routes, ok := shard.table[user_id]
	if ok {
		ok = false
		for i := 0; i < len(routes); i++ {
			if routes[i].session_id == session_id && routes[i].ws_id == ws_id && routes[i].sock_id == conn_obj.conn_id {
				ok = true
				shard.table[user_id][i].active = false
			}
		}
	}

	shard.Unlock()

}
