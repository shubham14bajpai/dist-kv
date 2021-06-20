package config

import (
	"fmt"
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

type Shard struct {
	Name    string
	Idx     int
	Address string
}

type Config struct {
	Shards []Shard
}

type Shards struct {
	Count   int
	CurrIdx int
	Addrs   map[int]string
}

func ParseFile(fileName string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(fileName, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func ParseShards(shards []Shard, currShardName string) (*Shards, error) {
	var shardCount = len(shards)
	var shardIdx int = -1
	var addrs = make(map[int]string)
	for _, s := range shards {

		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index %q found", s.Idx)
		}

		addrs[s.Idx] = s.Address
		if s.Name == currShardName {
			shardIdx = s.Idx
		}
	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard index %q not found", i)
		}
	}

	if shardIdx < 0 {
		return nil, fmt.Errorf("shard %s not found", currShardName)
	}

	return &Shards{
		Count:   shardCount,
		CurrIdx: shardIdx,
		Addrs:   addrs,
	}, nil
}

func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
