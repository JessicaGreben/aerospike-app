# aerospike app

Test app that writes data to an Aerospike db.

### Run app
Set env vars:
```
export AEROSPIKE_CLOUD_HOSTNAME=239f1a7c-0183-4de3-8407-172b833cadb9.aerospikedev.com
export AEROSPIKE_CLOUD_API_KEY_ID=962c9a9589d4455fab8b2c6d42e21f0a
export AEROSPIKE_CLOUD_API_KEY_SECRET=2f5466ae65c0479aa50cc21622873018
```

Build binary
```
$ go install
```

Usage:
```
$ aerospike-app -help
```

Run app once which will read/write a single record to the database to check for connectivity:
```
$ aerospike-app
```

Run app on a loop forever, sleeping every 15s between read/writes to the database. Each loop writes a single record to the db at the same key:
```
$ aerospike-app -forever
```

Run app to seed a database with fake data, the default settings runs 10 threads. Each thread creates 1000 rows each with 1024 bytes of data for a total of 10 mb of data.
```
$ aerospike-app -seed-data
```

Increase how many threads are used to seed the database with more fake data. Increasing concurrency to 100 will seed the database with 100MB of data:
```
$ aerospike-app -seed-data -concurrency 100
```
