package main

import (
	"fmt"
)

func cache_op(conn_obj *tcp_conn_object, cache_op_chan chan packaged_op, cache_op_util_chan chan []byte) {

	// OPEN/CREATE FILE FROM conn_obj socket type and ID

	var payload packaged_op
	var b []byte

	for {
		select {
		case payload = <-cache_op_chan:

			_ = payload // WRITE PAYLOAD TO DISK/SSD preferably ssd

		case b = <-cache_op_util_chan:
			switch b[0] {
			case 0: // close
				return
			case 1: // reset

				fmt.Println("Amount sent to socket:")

			case 2: // deploy reroute strategies
			}
		}
	}
}

//DEPRACATED (in memory cache)
// func cache_op(cache_op_chan chan packaged_op, cache_op_util_chan chan byte) {

// 	var op byte
// 	var payload packaged_op

// 	cur_cached_len := 10000
// 	cached := make([]packaged_op, cur_cached_len)
// 	cur_cached_head := 0
// 	cur_cached_tail := 0

// 	cached_leeway := 2
// 	cur_sequence_size := 0
// 	prev_sequence_sizes := make([]int, cached_leeway, cached_leeway)

// 	for {
// 		select {
// 		case payload = <-cache_op_chan:

// 			if cached[cur_cached_head] != nil {
// 				//fmt.Println(len(cached), cur_cached_head, cur_cached_tail)
// 				new_cached := make([]packaged_op, cur_cached_len*2)
// 				copy(new_cached[:cur_cached_head], cached[:cur_cached_head])
// 				copy(new_cached[cur_cached_len+cur_cached_head:], cached[cur_cached_head:])
// 				cur_cached_tail = cur_cached_len + cur_cached_head
// 				cur_cached_len *= 2
// 				cached = new_cached
// 			}

// 			cached[cur_cached_head] = payload

// 			if cur_cached_head == cur_cached_len-1 {
// 				cur_cached_head = 0
// 			} else {
// 				cur_cached_head++
// 			}
// 			cur_sequence_size++
// 		case op = <-cache_op_util_chan:
// 			switch op {
// 			case 0: // close
// 				return
// 			case 1: // reset
// 				if cached_leeway == 1 {

// 					for i := 0; i < prev_sequence_sizes[0]; i++ {
// 						cached[cur_cached_tail] = nil
// 						if cur_cached_tail == cur_cached_len-1 {
// 							cur_cached_tail = 0
// 						} else {
// 							cur_cached_tail++
// 						}
// 					}

// 					//fmt.Println(cur_cached_head, cur_cached_tail)

// 					copy(prev_sequence_sizes[:len(prev_sequence_sizes)-1], prev_sequence_sizes[1:])
// 					prev_sequence_sizes[len(prev_sequence_sizes)-1] = cur_sequence_size
// 				} else {
// 					prev_sequence_sizes[len(prev_sequence_sizes)-cached_leeway] = cur_sequence_size
// 					cached_leeway--
// 				}
// 				cur_sequence_size = 0
// 			case 2: // deploy reroute strategies
// 			}
// 		}
// 	}
// }
