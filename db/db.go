package db

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/button-tech/utils-node-tool/db/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	hosts    = "localhost:27017"
	database = "endpoints"
	// username   = ""
	// password   = ""
	collection = "addresses"
)

func GetEndpoint(currency string) (string, error) {
	session, err := mgo.Dial(hosts)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var addrs schema.Endpoints

	c := session.DB(database).C(collection)

	rand.Seed(time.Now().UnixNano())

	err = c.Find(bson.M{"currency": currency}).One(&addrs)
	if err != nil {
		return "", err
	}

	result := addrs.Addresses[rand.Intn(len(addrs.Addresses))]

	fmt.Println(result)

	return result, nil
}
