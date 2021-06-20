package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/shubham14bajpai/dist-kv/config"
	"github.com/shubham14bajpai/dist-kv/db"
)

type Server struct {
	db     *db.Database
	shards *config.Shards
}

func NewServer(db *db.Database, s *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: s,
	}
}

func (s *Server) redirect(shard int, rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(rw, "redirecting from shard %d to shard %d\n",
		s.shards.CurrIdx, shard)

	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
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

	shard := s.shards.Index(key)
	if shard != s.shards.CurrIdx {
		s.redirect(shard, rw, r)
		return
	}

	value, err := s.db.GetKey(key)
	fmt.Fprintf(rw, "shard = %d curr shard = %d Value = %q, error = %v\n",
		shard, s.shards.CurrIdx, value, err)
}

func (s *Server) SetHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.shards.Index(key)
	if shard != s.shards.CurrIdx {
		s.redirect(shard, rw, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(rw, "shard = %d curr shard = %d error = %v\n",
		shard, s.shards.CurrIdx, err)
}

func (s *Server) DeleteExtraKeysHandler(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(rw, "error cleaning extra keys: %v",
		s.db.DeleteExtraKeys(func(key string) bool {
			return s.shards.Index(key) != s.shards.CurrIdx
		}))

}
