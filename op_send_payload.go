package main

func send_payload(origin_session_id string, origin_id string, dest_id string, payload []byte) {

	var origin_user_routing_info user_routing_info
	var dest_user_routing_info user_routing_info
	var origin_ok bool
	var dest_ok bool

	origin_shard := routing_table.get_shard(origin_id)
	dest_shard := routing_table.get_shard(dest_id)

	if origin_shard == dest_shard {
		origin_shard.RLock()

		origin_user_routing_info, origin_ok = origin_shard.table[origin_id]
		dest_user_routing_info, dest_ok = origin_shard.table[dest_id]

		origin_shard.RUnlock()
	} else {
		origin_shard.RLock()

		origin_user_routing_info, origin_ok = origin_shard.table[origin_id]

		origin_shard.RUnlock()
		dest_shard.RLock()

		dest_user_routing_info, dest_ok = dest_shard.table[dest_id]

		dest_shard.RUnlock()
	}

	if origin_ok {
		for _, route := range origin_user_routing_info {
			if route.session_id != origin_session_id {
				if route.active {
					route_message_to_ws(route.sock_id, route.ws_id, 100, nil, payload)
				}
			}
		}
	}

	if dest_ok {
		for _, route := range dest_user_routing_info {
			if route.active {
				route_message_to_ws(route.sock_id, route.ws_id, 100, nil, payload)
			}
		}
	}

}
