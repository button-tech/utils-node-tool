# UTILS-NODE-TOOL

### Run in Docker container

Database for blockchain nodes addresses - MongoDB 
 
 - Build
```
# docker build --build-arg DIR=dir_with_main.go -t name_of_image .
```
- Run examples


```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection -e MAIN_API=url_of_main_api -e BLOCKCHAIN=eth-based(eth/etc)/utxo-based(btc,ltc,bch) name_of_image
```
Other

```
# docker run -p 8080:8080 -e HOST=db_host -e DB=database -e USER=db_user -e PASS=db_password -e COLLECTION=db_collection name_of_image
```
