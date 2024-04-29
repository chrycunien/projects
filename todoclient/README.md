# README

## Integration Test

### Build Constraints
```go
// After 1.17
//go:build

// Before 1.17
// +build
```

### Test
```bash
go test -v ./cmd -tags integration
go test ./cmd -tags integration -count=1
# normal
go test -v ./cmd
```