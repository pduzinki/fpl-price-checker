## About
*fpl-price-checker* is an app that processes data from Fantasy Premier League API, and generates daily player price changes reports. It can be used as a cli app, which saves data on *"/home/{current_user}/fpc"*, or as AWS Lambdas, which saves data on S3 Bucket specified with env variables, read by *config* package. For a succesful report generation for a given day, it needs players data from the previous day, and current day.

## Prerequisites
* Go (1.19 or higher) https://golang.org/doc/install
* Mage https://github.com/magefile/mage

## Building
```sh
# for cli app
mage cli

# for aws lambdas
mage lambdas
```

## Testing
```sh
go test ./...
```

## Running CLI
```sh
# fetches and saves current data on players available in FPL
./fpc fetch-fpl-data

# generates price changes report, based on FPL data obtained via fetch-fpl-data command
# to succeed, it needs two data points for comparison (i.e. FPL players data from the day before, and the day you run the command)
./fpc generate-report

# prints latest price changes report
./fpc get-report

# starts simple web server, with endpoints to get price changes report:
# /latest returns last generated report
# /:date returns report generated for given date (date format required "2006-01-02")
./fpc start-server
```

## License
*fpl-price-checker* is licensed under MIT License. See LICENSE file.