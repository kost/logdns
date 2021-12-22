# logdns
Simple DNS logging server

# Example

Simple 
```
./logdns
```

# Usage

```
Usage of ./logdns
  -listen string
        listen address (default "0.0.0.0")
  -logfile string
        Log File name
  -port int
        port to listen (default 53)
  -resolve string
        which domains to respond, e.g. service. (default ".")
  -return string
        what address to return (default "127.0.0.1")
  -ttl string
        Set a custom ttl for returned records (default "3600")
  -verbose
        be verbose
```

