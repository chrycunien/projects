# README

## SQLite
```bash
brew install sqlite3
go get github.com/mattn/go-sqlite3
```

```bash
# check
go env CGO_ENABLED
# set
go env -w CGO_ENABLED=1
```

## Testing
```bash
go test -v ./pomodoro
go test -v ./pomodoro -tags=inmemory
```