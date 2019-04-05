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
	hosts      = ""
	database   = ""
	username   = ""
	password   = ""
	collection = ""
)

var info = mgo.DialInfo{
	Addrs:[]string{hosts},
	Database:database,
	Username:username,
	Password:password,
}

func GetEndpoint(currency string) (string, error) {
	session, err := mgo.DialWithInfo(&info)
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

func GetAll() ([]schema.Endpoints, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var addresses []schema.Endpoints

	c := session.DB(database).C(collection)
	err = c.Find(nil).All(&addresses)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func DeleteAddress(currency, address string) (bool, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return false, err
	}
	defer session.Close()

	var entry schema.Endpoints

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"currency": currency}).One(&entry)
	if err != nil {
		return false, err
	}

	newAddresses := removeSliceElement(entry.Addresses, address)

	err = c.Update(bson.M{"currency": currency}, bson.M{"$set": bson.M{"addresses": newAddresses}})
	if err != nil {
		return false, err
	}

	return true, nil
}

func removeSliceElement(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
