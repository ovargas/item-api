# item-api


A simple API created as a PoC to play a bit with gRCP and protocol buffer

Build with Makefile

```bash
make
```

Build using `go`

```bash
go build -o bin/item-app cmd/main.go
```

Execute the service

```bash
./bin/item-app
```

Display command line options

```bash
./bin/item-app --help
```

```text
Usage of ./item-app:
  -port int
        The server port (default 10001)
  -storage_address string
        The storage server address in the format of host:port (default "localhost:10000")
```
