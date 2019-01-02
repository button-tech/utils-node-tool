# UTILS-NODE-TOOL

### Now available for: 
* Ethereum
* Ethereum Classic

### API:

### GET 
* "*/balance/:address" - return balance of account in crypto for specific nodeget balance of address
* "*/transactionFee" - return Amount of crypto that you need to send a transaction

* "*/gasPrice"

* "/eth/tokenBalance/:token/:address"



### Run on localhost:
```
# ETH_NODE=ADDRESS_OF_NODE go run eth/main.go
```
### Run in Docker container
```
# docker build -f Dockerfile.eth -t name .
# docker run -p 8080:8080 -e ETH_NODE=ADDRESS_OF_NODE name
```