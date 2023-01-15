## About
*fpl-price-checker* is an app that can be used to process data from Fantasy Premier League API, and generate reports on daily player price changes. 

## Prerequisites
* Go (1.19 or higher) https://golang.org/doc/install
* Mage https://github.com/magefile/mage

## Building
```sh
#with mage
mage build

#with go build
go build -o fpc ./cmd
```

## Testing
```sh
go test ./...
```

## Running
```sh
# fetches and saves current data on players available in FPL
./fpc fetch-fpl-data

# generates price changes report, based on FPL data obtained via fetch-fpl-data command
# to succeed, it needs two data points for comparison (i.e. FPL players data from the day before, and the day you run the command)
./fpc generate-report

# starts simple web server, with endpoints to get price changes report:
# /latest returns last generated report
# /:date returns report generated for given date (date format required "2006-01-02")
./fpc start-server
```

## License
*fpl-price-checker* is licensed under MIT License. See LICENSE file.