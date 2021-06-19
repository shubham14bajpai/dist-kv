package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/shubham14bajpai/dist-kv/config"
	"github.com/shubham14bajpai/dist-kv/db"
	"github.com/shubham14bajpai/dist-kv/web"
)

var (
	dbLocation = flag.String("db-location", "", "path for bolt db location")
	httpAddr   = flag.String("http-addr", ":8080", "http host and port")
	configFile = flag.String("config-file", "sharding.toml", "config file for static sharding")
	shard      = flag.String("shard", "", "the name for the shard for the data")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("must provide db location")
	}

	if *shard == "" {
		log.Fatal("must provide shard name")
	}
}

func main() {

	parseFlags()
	var c config.Config
	_, err := toml.DecodeFile(*configFile, &c)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)
	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address
		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("shard %s not found", *shard)
	}

	log.Printf("Shard count: %d and shard index: %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer close()

	srv := web.NewServer(db, shardCount, shardIdx, addrs)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}
