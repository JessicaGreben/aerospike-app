package main

import (
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aerospike/aerospike-client-go/v6"
)

const (
	namespace = "aerospike_cloud"
)

var (
	forever     bool
	concurrency int
	seedData    bool
)

func main() {
	flag.BoolVar(&forever, "forever", false, "loops forever reading/writing a single record to the database to ensure connectivity. Defaults to false.")
	flag.IntVar(&concurrency, "concurrency", 10, "how many concurrent goroutines to use when writing data.")
	flag.BoolVar(&seedData, "seed-data", false, "seed the database with fake data. Defaults to false.")
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

	t := newTester(client)
	if err := t.run(); err != nil {
		log.Fatal(err)
	}
}

type tester struct {
	asClient    *aerospike.ProxyClient
	readPolicy  *aerospike.BasePolicy
	writePolicy *aerospike.WritePolicy
}

func newTester(client *aerospike.ProxyClient) *tester {
	writePolicy := aerospike.NewWritePolicy(0, 0)
	writePolicy.TotalTimeout = 5 * time.Second
	readPolicy := aerospike.NewPolicy()
	readPolicy.TotalTimeout = 5 * time.Second
	return &tester{
		asClient:    client,
		readPolicy:  readPolicy,
		writePolicy: writePolicy,
	}
}

func (t *tester) run() error {
	start := time.Now().UTC()
	defer log.Printf("Total time of program: %v", time.Since(start))
	if seedData {
		return t.seedData()
	}
	return t.testReadWriteConnectivity()
}

func (t *tester) seedData() error {
	log.Println("seeding data")
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		time.Sleep(1 * time.Second) // wait a little bit to start next goroutine so we dont bombard the db all at once
		go func(id int) {
			t.createFakeData1kRecords(id, "users")
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}

func (t *tester) createFakeData1kRecords(id int, setName string) {
	const recordCount = 1000
	setNamePerWorker := fmt.Sprintf("%d%s", id, setName)
	log.Printf("threadID=%d creating %d records of data in set=%s\n", id, recordCount, setNamePerWorker)

	binData := make([]byte, 1024)
	if _, err := rand.Read(binData); err != nil {
		log.Println("rand.Read err: ", err)
	}

	start := time.Now().UTC()

	for i := 0; i < recordCount; i++ {
		key := fmt.Sprintf("%d-%d-keyname", i, id)
		asKey, err := aerospike.NewKey(namespace, setNamePerWorker, key)
		if err != nil {
			log.Println("newKey err: ", err)
		}
		bin := aerospike.NewBin("b", binData)
		if err := t.asClient.PutBins(t.writePolicy, asKey, bin); err != nil {
			log.Println("PutBins err: ", err)
		}
	}

	fmt.Printf("threadID=%d  Time to create %d records: %s\n", id, recordCount, time.Since(start))
}

func (t *tester) testReadWriteConnectivity() error {
	key, err := aerospike.NewKey(namespace, "foo", "bar")
	if err != nil {
		log.Println("newKey err: ", err)
	}

	for {
		t.readWrite(key)
		if !forever {
			return nil
		}
		time.Sleep(15 * time.Second)
	}
}

func (t *tester) readWrite(key *aerospike.Key) {
	bin1 := aerospike.NewBin("b", fmt.Sprintf("data"+time.Now().String()))
	bin2 := aerospike.NewBin("secondbin", "data2")

	err := t.asClient.PutBins(t.writePolicy, key, bin1, bin2)
	if err != nil {
		log.Println("Put err: ", err)
	} else {
		log.Println("Succesfully wrote record")
	}

	record, err := t.asClient.Get(t.readPolicy, key)
	if err != nil {
		log.Println("Get err: ", err)
		return
	}

	log.Printf("Reading record: key: %s, bin map:%v\n", record.Key, record.Bins)
}
