# aerospike app

Test app that writes a single record to an Aerospike db then reads it back.

### Run app
Set env vars:
```
export AEROSPIKE_CLOUD_HOSTNAME=
export AEROSPIKE_CLOUD_API_KEY_ID=
export AEROSPIKE_CLOUD_API_KEY_SECRET=
```

Run app once:
```
$ go run main.go
```

Run app on a loop forever, sleeping every 15s between read/writes to the database:
```
$ go run main.go --forever
```
