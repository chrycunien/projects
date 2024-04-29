# README

## Dependency
```bash
go mod edit -require=todolist@v0.0.0
go mod edit -replace=todolist=../todolist
# check
go list -m all
```