package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aerospike/aerospike-client-go/v6"
)

func main() {
	forever := flag.Bool("forever", false, "loop forever")
	flag.Parse()

	aerospikeHost := os.Getenv("AEROSPIKE_CLOUD_HOSTNAME")
	if aerospikeHost == "" {
		log.Fatal("required env var AEROSPIKE_CLOUD_HOSTNAME is not set.")
	}
	aerospikeKey := os.Getenv("AEROSPIKE_CLOUD_API_KEY_ID")
	if aerospikeKey == "" {
		log.Fatal("required env var AEROSPIKE_CLOUD_API_KEY_ID is not set.")
	}
	aerospikeSecret := os.Getenv("AEROSPIKE_CLOUD_API_KEY_SECRET")
	if aerospikeSecret == "" {
		log.Fatal("required env var AEROSPIKE_CLOUD_API_KEY_SECRET is not set.")
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

	key, err := aerospike.NewKey("aerospike_cloud", "foo", "bar")
	if err != nil {
		log.Println("0: ", err)
	}

	if *forever {
		for {
			do(client, key)
			time.Sleep(15 * time.Second)
		}
	}
	do(client, key)
}

func do(client *aerospike.ProxyClient, key *aerospike.Key) {
	writePolicy := aerospike.NewWritePolicy(0, 0)
	writePolicy.TotalTimeout = 5 * time.Second
	bin1 := aerospike.NewBin("firstbin", fmt.Sprintf("data"+time.Now().String()))
	bin2 := aerospike.NewBin("secondbin", "data2")

	err := client.PutBins(writePolicy, key, bin1, bin2)
	if err != nil {
		log.Println("1: ", err)
	}
	log.Println("Succesfully wrote record")

	readPolicy := aerospike.NewPolicy()
	readPolicy.TotalTimeout = 5 * time.Second
	record, err := client.Get(readPolicy, key)
	if err != nil {
		log.Println("2: ", err)
	}

	log.Printf("Reading record: key: %s, bin map:%v\n", record.Key, record.Bins)
}
