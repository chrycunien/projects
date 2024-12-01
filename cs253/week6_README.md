# Week6

## Seventeen
```bash
cd week6/seventeen

# Option 1:
echo "main.DataStorageManager" | make
echo "main.StopwordManager" | make
echo "main.WordFrequencyManager" | make

# Option 2:
go run ./ ../../pride-and-prejudice.txt
```

## Twenty
```bash
# change the config.ini for different options
cd week6/twentytwo

# Option 1:
make run

# Option 2:
go build -buildmode=plugin -o ./app1/app1.so ./app1/app1.go
go build -buildmode=plugin -o ./app2/app2.so ./app2/app2.go
go run ./ ../../pride-and-prejudice.txt
```

## Result
The result file is put in `pride-and-prejudice-result.txt` for easy comparison.