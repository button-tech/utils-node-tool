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



### ETH on localhost:
```
# ETH_NODE=ADDRESS_OF_NODE go run eth/main.go
```
### ETH in docker container
```
# docker build -f Dockerfile.eth -t name .
# docker run -p 8080:8080 -e ETH_NODE=ADDRESS_OF_NODE name
```

### ETC on localhost:
```
# ETC_NODE=ADDRESS_OF_NODE go run etc/main.go
```
### ETH in docker container
```
# docker build -f Dockerfile.etc -t name .
# docker run -p 8080:8080 -e ETC_NODE=ADDRESS_OF_NODE name
```