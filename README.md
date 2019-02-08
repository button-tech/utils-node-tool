# UTILS-NODE-TOOL

### API:

### GET 
For all:
* "*/balance/:address" - return balance of account in crypto for specific nodeget balance of address

For eth,etc,btc,ltc,bch:
* "*/transactionFee" - return Amount of crypto that you need to send a transaction

For eth,etc:
* "*/gasPrice"

For eth:
* "/eth/tokenBalance/:token/:address"

For btc,ltc,bch:
* "*/utxo/:address"

### POST:
For all:

* "*/balances" -> send "addressesArray":["address","address"]



### Run on localhost:
```
# ETH_NODE=ADDRESS_OF_NODE go run eth/main.go
```
### Run in Docker container
```
# docker build -f Dockerfile.eth -t name .
# docker run -p 8080:8080 -e ETH_NODE=ADDRESS_OF_NODE name
```
