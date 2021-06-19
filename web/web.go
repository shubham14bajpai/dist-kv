package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/shubham14bajpai/dist-kv/db"
)

type Server struct {
	db         *db.Database
	shardCount int
	shardIdx   int
	addrs      map[int]string
}

func NewServer(db *db.Database, shardCount, shardIdx int, addrs map[int]string) *Server {
	return &Server{
		db:         db,
		shardCount: shardCount,
		shardIdx:   shardIdx,
		addrs:      addrs,
	}
}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(shard int, rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "redirecting from shard %d to shard %d\n", s.shardIdx, shard)
	url := "http://" + s.addrs[shard] + r.RequestURI
	resp, err := http.Get(url)
	if err != nil {
		rw.WriteHeader(500)
		fmt.Fprintf(rw, "failed to redirect request: %v", err)
		return
	}

	defer resp.Body.Close()
	io.Copy(rw, resp.Body)
}

func (s *Server) GetHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.getShard(key)
	if shard != s.shardIdx {
		s.redirect(shard, rw, r)
		return
	}
	value, err := s.db.GetKey(key)
	fmt.Fprintf(rw, "shard = %d curr shard = %d Value = %q, error = %v\n", shard, s.shardIdx, value, err)
}

func (s *Server) SetHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	shard := s.getShard(key)
	if shard != s.shardIdx {
		s.redirect(shard, rw, r)
		return
	}
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(rw, "shard = %d curr shard = %d error = %v\n", shard, s.shardIdx, err)
}
