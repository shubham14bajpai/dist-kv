package main

import (
	"flag"
	"log"
	"net/http"

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

	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("failed to parse shards: %v", err)
	}

	log.Printf("Shard count: %d and shard index: %d",
		shards.Count, shards.CurrIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer close()

	srv := web.NewServer(db, shards)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeysHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}
