package main

import (
	"errors"
)

func handle_socket_op(conn_obj *tcp_conn_object, b []byte, op byte) {

	var err error = nil

	switch op {

	case 20:
		param, _, err := byte_to_params_and_payload(b, 3, false)
		if err != nil {
			break
		}
		connect_session(conn_obj, param[0], param[1], param[2])
	case 21:
		param, _, err := byte_to_params_and_payload(b, 4, false)
		if err != nil {
			break
		}
		recover_session(conn_obj, param[0], param[1], param[2], param[3])
	case 22:
		param, _, err := byte_to_params_and_payload(b, 3, false)
		if err != nil {
			break
		}
		remove_session(conn_obj, param[0], param[1], param[2])
	case 30:
		param, payload, err := byte_to_params_and_payload(b, 3, true)
		if err != nil {
			break
		}
		send_prev_state(param[0], param[1], param[2], payload)
	default:
		err = errors.New("default switch triggered")
	}

	if err != nil {
		log_err("socket op:" + string(op) + " error:" + err.Error())
	}

}

func handle_api_op(conn_obj *tcp_conn_object, b []byte, op byte) {

	var err error = nil

	switch op {
	case 100:
		param, payload, err := byte_to_params_and_payload(b, 3, true)
		if err != nil {
			break
		}
		send_payload(param[0], param[1], param[2], payload)
	default:
		err = errors.New("default switch triggered")
	}

	if err != nil {
		log_err("api op:" + string(op) + " error:" + err.Error())
	}

}
