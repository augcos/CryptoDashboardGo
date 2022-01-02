# CryptoDashboardGo
## Introduction
This is a Go implemetation of a dashboard to track your cryptocurrency investments in Binance. You can find the original version developed in Python [here](https://github.com/augcos/CryptoDashboardPy). This project was developed using Go v1.17.5 for Linux systems.

## How to use
First, clone the repository to your local system:
```
git clone https://github.com/augcos/CryptoDashboardGo
```
Before running the dashboard, you will need to add a CSV file with your trading history to the cloned directory. You can download it from your Binance profile. Make sure it has the long format structure provided by Binance. Then, run the code using the command.
```
go run dashboard.go --csv [filename.csv] --refresh [refresh time in seconds] --order [display order]
```
You can also compile the code yourself and then run the binary:
```
go build dashboard.go
./dashboard --csv [filename.csv] --refresh [refresh time in seconds] --order [display order]
```
Remember to use the flags to overwrite the default settings (csv: myOrders.csv, refresh limit: 60 seconds, order: "Invested BTC").