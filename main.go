package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ProSellers/go-honeyapi/api/wb"
	"github.com/ProSellers/go-honeyapi/internal/controllers"
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/utils"
	"github.com/ProSellers/go-honeyapi/utils/cfg"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/sirupsen/logrus"
)

func main() {
	// wb.GetBrandIdAndUseFirst()ev
	cfg.Load()

	_, err := database.Init()
	if err != nil {
		panic(err)
	}

	// prox := controllers.NewProxyController()
	// prox.AddProxy(&controllers.Proxy{Type: 2, IP: "10.124.0.77", Port: 1080}, &controllers.Proxy{Type: 2, IP: "10.124.0.152", Port: 1080})

	// go StartScanner(db, prox)

	err = startServer()
	if err != nil {
		logrus.Fatalln(err)
	}
}

func StartScanner(db *database.Db, prox *controllers.ProxyController) {
	for {
		var brandcount int64

		fmt.Println("starting scanner")

		var brands = make(chan wb.Brand, 100000000)

		utils.UseAlphabet(func(letter string) error {
			api := wb.New()
			// ppx, err := prox.GetProxy()
			// if err != nil {
			// 	if !errors.Is(err, controllers.ErrNoProxies) {
			// 		return err
			// 	}
			// } else {
			// 	api.SetProxy(ppx)
			// }

			list, err := api.GetBrandList(letter)
			if err != nil {
				return err
			}

			fmt.Println(letter + " | Got " + strconv.Itoa(len(list.Value.BrandsList)) + " brands")

			atomic.AddInt64(&brandcount, int64(len(list.Value.BrandsList)))

			for _, el := range list.Value.BrandsList {
				brands <- el
			}
			return nil
		})

		fmt.Println(len(brands))

		var wg sync.WaitGroup

		var products = make(chan []wb.Products, 100000000)

		var productcount int64
		var scanned int64

		go func() {
			for {
				fmt.Println(strconv.FormatInt(scanned, 10) + "/" + strconv.FormatInt(brandcount, 10) + " | Products parsed: " + strconv.FormatInt(productcount, 10))
				time.Sleep(time.Second)
			}
		}()

		close(brands)

		var bq = goconcurrentqueue.NewFIFO()

		for el := range brands {
			bq.Enqueue(el)
		}

		for i := 0; i <= 120; i++ {
			wg.Add(1)
			go func() {
				for {
					bbq, err := bq.Dequeue()
					if err != nil {
						break
					}

					brand := bbq.(wb.Brand)

					api := wb.New()

					ppx, err := prox.GetProxy()
					if err != nil {
						if !errors.Is(err, controllers.ErrNoProxies) {
							bq.Enqueue(brand)
							continue
						}
					} else {
						api.SetProxy(ppx)
					}

					info, err := api.GetBrand(brand.ID)
					if err != nil {
						continue
					}

					products <- info.Data.Products

					go atomic.AddInt64(&scanned, 1)
					go atomic.AddInt64(&productcount, int64(len(info.Data.Products)))
				}
				wg.Done()
			}()
		}

		wg.Wait()

		close(products)

		fmt.Println(len(products))

		time.Sleep(time.Hour * 6)
	}
}
