package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/anvie/port-scanner"
	"github.com/button-tech/utils-node-tool/shared/db"
)

func main() {

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	for {

		entrys, err := db.GetAll()
		if err != nil {
			fmt.Println(err)
		}

		for _, entry := range entrys {

			for _, address := range entry.Addresses {

				ip := re.FindString(address)

				ps := portscanner.NewPortScanner(ip, 5*time.Second, 1)

				isAlive := ps.IsOpen(entry.Port)

				if !isAlive {
					time.Sleep(time.Second * 10)
					secondCheck := ps.IsOpen(entry.Port)
					if !secondCheck {
						isDel, err := db.AddToStoppedList(entry.Currency, address)
						if err != nil {
							log.Println(err)
						}
						if !isDel {
							panic("Cant add")
						} else {
							fmt.Printf("Add to stopped list: %s", address)
						}
					}
				}
			}
		}

		fmt.Println("\nAll nodes checked!")
		time.Sleep(time.Second * 5)
	}
}
