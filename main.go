package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/aerospike/aerospike-client-go/v6"
)

func main() {
	aerospikeHost := os.Getenv("AEROSPIKE_CLOUD_HOST")
	if aerospikeHost == "" {
		log.Fatal("env var AEROSPIKE_CLOUD_HOST not set.")
	}
	aerospikeKey := os.Getenv("AEROSPIKE_CLOUD_KEY")
	if aerospikeKey == "" {
		log.Fatal("env var AEROSPIKE_CLOUD_KEY not set.")
	}
	aerospikeSecret := os.Getenv("AEROSPIKE_CLOUD_SECRET")
	if aerospikeSecret == "" {
		log.Fatal("env var AEROSPIKE_CLOUD_SECRET not set.")
	}

	clientPolicy := aerospike.NewClientPolicy()
	clientPolicy.User = aerospikeKey
	clientPolicy.Password = aerospikeSecret
	clientPolicy.TlsConfig = &tls.Config{}

	client, err := aerospike.NewProxyClient(
		clientPolicy,
		aerospike.NewHost(aerospikeHost, 4000),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	writePolicy := aerospike.NewWritePolicy(0, 0)
	writePolicy.TotalTimeout = 5 * time.Second

	key, err := aerospike.NewKey("aerospike_cloud", "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	bin := aerospike.NewBin("firstbin", "data")

	err = client.PutBins(writePolicy, key, bin)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Succesfully wrote record")

	readPolicy := aerospike.NewPolicy()
	readPolicy.TotalTimeout = 5 * time.Second
	record, err := client.Get(readPolicy, key)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Reading record: key: %s, bin map:%v", record.Key, record.Bins)
}
