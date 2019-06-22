package main

import (
	"log"
	"regexp"
	"time"

	"errors"
	"github.com/anvie/port-scanner"
	"github.com/button-tech/utils-node-tool/shared/db"
	"golang.org/x/sync/errgroup"
)

func main() {

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	var g errgroup.Group

	for {

		entrys, err := db.GetAll()
		if err != nil {
			log.Println(err)
			break
		}

		for _, entry := range entrys {

			for _, address := range entry.Addresses {

				address := address

				g.Go(func() error {

					ip := re.FindString(address)

					ps := portscanner.NewPortScanner(ip, 5*time.Second, 1)

					isAlive := ps.IsOpen(entry.Port)

					if !isAlive {

						time.Sleep(time.Second * 10)

						secondCheck := ps.IsOpen(entry.Port)

						if !secondCheck {

							isDel, err := db.AddToStoppedList(entry.Currency, address)
							if err != nil {
								return err
							}
							if !isDel {
								return errors.New("Can't del!\n")
							} else {
								log.Printf("Add to stopped list: %s", address)
							}

						}
					}
					return nil
				})
			}
			if err := g.Wait(); err != nil {
				log.Println(err)
				time.Sleep(time.Second * 10)
			}
		}

		log.Println("\nAll nodes checked!")

		time.Sleep(time.Second * 10)
	}
}
