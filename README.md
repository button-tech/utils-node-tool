# UTILS-NODE-TOOL

### API:

### GET 

####For all(eth,etc,btc,ltc,bch,waves,stellar):
* "*/balance/:address" - return balance of account in crypto for specific nodeget balance of address

####For eth,etc,btc,ltc,bch:
* "*/transactionFee" - return Amount of crypto that you need to send a transaction

####For eth,etc:
* "*/gasPrice"

####For eth:
* "/eth/tokenBalance/:token/:address"

####For btc,ltc,bch:
* "*/utxo/:address"

### POST:

####For all:

* "*/balances" -> send "addressesArray":["address","address"]


### Run in Docker container
 - Build
```
# docker build -f --build-arg DIR=dir_with_main.go -t name_of_image .
```
- Run (env for MongoDB)

```
# docker run -p 8080:8080 -e HOST=host -e DB=database -e USER=username -e PASS=password -e COLLECTION=collection name_of_image
```
