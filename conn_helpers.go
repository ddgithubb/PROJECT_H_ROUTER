package main

import (
	"encoding/binary"
	"errors"
	"strings"
)

type packaged_op []byte

const DEFAULT_SPLIT_CHAR = ":"

var DEFAULT_DELIM byte = DEFAULT_SPLIT_CHAR[0]

// Converts binary data to string parameters and payload
func byte_to_params_and_payload(b []byte, n_param byte, payload_exists bool) (param []string, payload []byte, err error) {

	param = make([]string, n_param)
	err = nil

	if n_param == 0 {
		if payload_exists {
			payload = b
		} else if len(b) > 1 || (len(b) == 1 && b[0] != DEFAULT_DELIM) {
			err = errors.New("unexpected payload")
		}
		return
	}

	var n byte = 0
	last_delim_i := 0
	for i := 0; i < len(b); i++ {
		if n == n_param-1 && !payload_exists {
			param[n] = string(b[last_delim_i:])
			return
		}
		if b[i] == DEFAULT_DELIM {
			param[n] = string(b[last_delim_i:i])
			last_delim_i = i + 1
			n++
			if n == n_param {
				if payload_exists {
					payload = b[last_delim_i:]
				} else {
					if i == len(b)-1 {
						return
					}
					err = errors.New("expected less parameters")
				}
				return
			}
		}
	}

	err = errors.New("expected more parameters")
	return
}

// func params_and_payload_to_bytes(params []string, payload []byte) (b []byte) {

// 	if params != nil {
// 		b_param := []byte(strings.Join(params, DEFAULT_SPLIT_CHAR))
// 		if payload != nil {
// 			b = make([]byte, len(b_param) + len(payload) + 1)
// 			b[len(b_param)] = DEFAULT_DELIM
// 			_ = copy(b[:len(b_param)], b_param)
// 			_ = copy(b[len(b_param) + 1:], payload)
// 		} else {
// 			b = b_param
// 		}
// 	} else if payload != nil {
// 		b = payload
// 	}

// 	return
// }

func package_op(op byte, params []string, payload []byte) packaged_op {

	b_param := []byte(strings.Join(params, DEFAULT_SPLIT_CHAR))

	if params != nil && payload != nil {
		b_param = append(b_param, DEFAULT_DELIM)
	}
	size := len(b_param) + len(payload)

	packaged := make([]byte, 5+size)
	packaged[0] = op
	binary.BigEndian.PutUint32(packaged[1:5], uint32(size))
	copy(packaged[5:5+len(b_param)], b_param)
	copy(packaged[5+len(b_param):], payload)

	return packaged
}

func package_op_router_wrapper(op byte, params []string, underlying_op byte, underlying_params []string, underlying_payload []byte) packaged_op {

	//fmt.Println("PREPARSE:", op, params, underlying_op, underlying_params, underlying_payload)

	b_param := []byte(strings.Join(params, DEFAULT_SPLIT_CHAR) + DEFAULT_SPLIT_CHAR)
	b_underlying_param := []byte(strings.Join(underlying_params, DEFAULT_SPLIT_CHAR))

	if underlying_params != nil && underlying_payload != nil {
		b_underlying_param = append(b_underlying_param, DEFAULT_DELIM)
	}

	size := len(b_param) + 1 + len(b_underlying_param) + len(underlying_payload)

	packaged := make([]byte, 5+size)
	packaged[0] = op
	binary.BigEndian.PutUint32(packaged[1:5], uint32(size))
	copy(packaged[5:5+len(b_param)], b_param)
	packaged[5+len(b_param)] = underlying_op
	copy(packaged[5+len(b_param)+1:5+len(b_param)+1+len(b_underlying_param)], b_underlying_param)
	copy(packaged[5+len(b_param)+1+len(b_underlying_param):], underlying_payload)

	return packaged
}
