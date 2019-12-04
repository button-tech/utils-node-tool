package db

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/button-tech/utils-node-tool/db/schema"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

var (
	host       = os.Getenv("HOST")
	database   = os.Getenv("DB")
	username   = os.Getenv("USER")
	password   = os.Getenv("PASS")
	collection = os.Getenv("COLLECTION")
)

var info = mgo.DialInfo{
	Addrs:    []string{host},
	Database: database,
	Username: username,
	Password: password,
}

func GetEndpoint(currency string) (string, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var addrs schema.EndpointsData

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

func GetAll() ([]schema.EndpointsData, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var addresses []schema.EndpointsData

	c := session.DB(database).C(collection)
	err = c.Find(nil).All(&addresses)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func AddToStoppedList(currency, address string) error {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return err
	}
	defer session.Close()

	var entry schema.EndpointsData

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"currency": currency}).One(&entry)
	if err != nil {
		return err
	}

	newAddresses := removeSliceElement(entry.Addresses, address)

	stoppedList := append(entry.Stopped, address)

	err = c.Update(bson.M{"currency": currency}, bson.M{"$set": bson.M{"stopped": stoppedList, "addresses": newAddresses}})
	if err != nil {
		return err
	}

	return nil
}

func GetEntry() (*schema.EndpointsData, error) {
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	var entry schema.EndpointsData

	c := session.DB(database).C(collection)

	err = c.Find(bson.M{"currency": os.Getenv("BLOCKCHAIN")}).One(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func removeSliceElement(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func DeleteEntry(currency, address string) error {
	err := AddToStoppedList(currency, address)
	if err != nil {
		return err
	}

	log.Printf("Add to stopped list %s: %s", currency, address)

	return nil
}
