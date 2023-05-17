# Redis Performance Test

This project is a rudimentary URL shortener to check how fast concurrent Redis calls (using go routines) are compared to sequential calls using Go.

## Prerequisites

- Redis installed (I used version `7.0.11`)
- Go installed (I used version `1.20.4`)

## Set Up and Perform Tests

**Clone repo**

```bash
git clone https://github.com/dominik-schmid/redis-performance-test.git
```

**Start Redis with config**

```bash
redis-server ./redis.conf`
```

**Install go modules**

```bash
go mod tidy
```

**Perform tests**

```bash
go run ./cmd/redis-performance-test.go
```

## Example Results

```bash
Connection to Redis successfully established
Looking up 200000 MD5 hashes sequentially took 7.83s
Looking up 200000 MD5 hashes using go routines took 2.36s
```

## Possible Improvements

- Perform concurrent tests with different configurations, i.e. stagger number of concurrent function calls to determine the performance sweet spot
- Improve output to make the comparison of the configurations easier
- Write tests for this project