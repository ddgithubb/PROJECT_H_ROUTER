package main

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/segmentio/fasthash/fnv1a"
)

type conn_types byte

const (
	SOCKET_TYPE = 0
	API_TYPE    = 1
	ROUTER_TYPE = 2
)

type tcp_conn_object struct {
	conn       *net.TCPConn
	conn_id    string
	conn_type  conn_types
	conn_group string
}

type unsafe_table struct {
	p unsafe.Pointer
	sync.Mutex
}

type session_routing_info struct {
	session_id string
	ws_id      string
	sock_id    string
	active     bool
}
type user_routing_info []*session_routing_info

type conc_routing_table struct {
	table map[string]user_routing_info
	sync.RWMutex
}
type conc_routing_table_shards []*conc_routing_table

type cache_writer struct {
	writer io.Writer
	sync.Mutex
}

type conn_route struct {
	conn_obj *tcp_conn_object

	cache *cache_writer
}
type conn_table map[string]*conn_route
type conn_table_ptr *conn_table

var socket_table *unsafe_table = new_conn_table()
var api_table *unsafe_table = new_conn_table()

var routing_table conc_routing_table_shards = func() conc_routing_table_shards {
	shards := make([]*conc_routing_table, CONCURRENCY)

	for i := 0; uint32(i) < CONCURRENCY; i++ {
		shards[i] = &conc_routing_table{table: make(map[string]user_routing_info)}
	}

	return shards
}()

func (crt conc_routing_table_shards) get_shard(id string) *conc_routing_table {
	return crt[fnv1a.HashString32(id)%CONCURRENCY]
}

func new_conn_table() *unsafe_table {
	var ptr unsafe.Pointer

	conn_table := make(conn_table)
	atomic.StorePointer(&ptr, unsafe.Pointer(&conn_table))

	return &unsafe_table{
		p: ptr,
	}
}
