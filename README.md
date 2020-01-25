# UTILS-NODE-TOOL

### Blockchain list: ETH based(ETH/ETC),UTXO based(BTC/BCH/LTC), WAVES, STELLAR, ZILLIQA, POA, TON(testnet)

### Run in Docker container

Database for blockchain nodes addresses - MongoDB 
 
 - Build
```
# docker build --build-arg DIR=dir_in_cmd -t name_of_image .
```
- Run examples

ETH/UTXO Based:
```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection -e BLOCKCHAIN=(eth/etc/poa/btc/ltc/bch) -e DSN=(sentry DSN) name_of_image
```

WAVES/XLM:

```
# docker run -p 8080:8080 name_of_image
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
DSN=...
```
2. Build and run:
```
# docker-compose build
# docker-compose up
```
