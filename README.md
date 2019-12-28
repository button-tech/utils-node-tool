# UTILS-NODE-TOOL

### Blockchain list: ETH based(ETH/ETC),UTXO based(BTC/BCH/LTC), WAVES, XLM, ZILLIQA, COSMOS, RIPPLE, TEZOS, ALGORAND, BNB, TON(testnet)

### Run in Docker container

Database for blockchain nodes addresses - MongoDB 
 
 - Build
```
# docker build --build-arg DIR=dir_in_cmd -t name_of_image .
```
- Run examples


```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection -e MAIN_API=url_of_main_api -e BLOCKCHAIN=eth-based(eth/etc)/utxo-based(btc,ltc,bch) name_of_image
```
Other

```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection name_of_image
```

### TON(testnet)

```
# docker build -f Dockerfile.ton -t name_of_image .
# docker run -p 8080:8080 -e WORKDIR=/app/ name_of_image
```

### Docker-compose
1. Create .env file:
```
HOST=...
DB=...
COLLECTION=...
USER=...
PASS=...
```
2. Build and run:
```
# docker-compose build
# docker-compose up
```
