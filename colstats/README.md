# README

## Requirements
Install graphiz for pprof web interface
```bash
brew install graphviz
```

## Profiling
```bash
go test -bench . -benchtime=10x -cpuprofile cpu00.pprof
go test -bench . -benchtime=10x -memprofile mem00.pprof
go test -bench . -benchtime=10x -benchmem | tee benchresults00m.txt
go get -u -v golang.org/x/tools/cmd/benchcmp
go install golang.org/x/tools/cmd/benchcmp
benchcmp benchresults00m.txt benchresults01m.txt

# Note: benchstat is to replace deprecated benchcmp
go get golang.org/x/perf/cmd/benchstat
go install golang.org/x/perf/cmd/benchstat
```

### pprof
```bash
go tool pprof cpu00.pprof 

# get stats
(pprof) top
# cumulative
(pprof) top -cum
# list func
(pprof) list csv2float
# web
(pprof) web
# quit
(pprof) quit
```

## Tracing
```bash
go test -bench . -benchtime=10x -trace trace01.out
```

### trace
```bash
go tool trace trace01.out
```