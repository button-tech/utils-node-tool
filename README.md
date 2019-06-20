# UTILS-NODE-TOOL

### Run in Docker container

Database for blockchain nodes addresses - MongoDB 
 
 - Build
```
# docker build -f --build-arg DIR=dir_with_main.go -t name_of_image .
```
- Run examples

ETH based(ETH/ETC)
```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection -e main-api=url_of_main_api -e blockchain=eth_based_blockchain(eth or etc)  name_of_image
```

UTXO based(BTC/LTC/BCH)
```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection -e main-api=url_of_main_api -e blockchain=utxo_based_blockchain(btc, bch or ltc)  name_of_image
```

Other

```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection name_of_image
```
