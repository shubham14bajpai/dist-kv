package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParseShards(t *testing.T) {
	type args struct {
		shards        []Shard
		currShardName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Shards
		wantErr bool
	}{
		{
			name: "test parsing of shards",
			args: args{
				shards: []Shard{
					{
						Name:    "north",
						Idx:     0,
						Address: "localhost:8080",
					},
					{
						Name:    "east",
						Idx:     1,
						Address: "localhost:8081",
					},
				},
				currShardName: "east",
			},
			want: &Shards{
				Count:   2,
				CurrIdx: 1,
				Addrs: map[int]string{
					0: "localhost:8080",
					1: "localhost:8081",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseShards(tt.args.shards, tt.args.currShardName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseShards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseShards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTempFile() (string, error) {
	contents := `[[shards]]
	name = "north"
	idx = 0
	address = "localhost:8080"
	`
	var name string

	f, err := ioutil.TempFile(os.TempDir(), "config.toml")
	if err != nil {
		return name, err
	}
	defer f.Close()

	if _, err := f.WriteString(contents); err != nil {
		return name, err
	}

	name = f.Name()
	return name, nil
}

func TestParseFile(t *testing.T) {
	tempFile, err := getTempFile()
	if err != nil {
		t.Fatalf("failed to generate temp file: %v", err)
		return
	}
	defer os.Remove(tempFile)
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "test config file parsing",
			args: args{
				fileName: tempFile,
			},
			want: Config{
				[]Shard{
					{
						Name:    "north",
						Idx:     0,
						Address: "localhost:8080",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
