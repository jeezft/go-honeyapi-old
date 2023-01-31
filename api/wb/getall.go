package wb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
)

type Good struct {
	Sort  int    `json:"sort"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

func GetBrandIdAndUseFirst() {
	a := New()

	brands, err := a.GetBrandList("a")
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	var queue = goconcurrentqueue.NewFIFO()

	for _, brand := range brands.Value.BrandsList {
		queue.Enqueue(brand.ID)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for {
				b, err := queue.Dequeue()
				if err != nil {
					break
				}

				brand := b.(int)

				info, err := a.GetBrand(brand)
				if err != nil {
					fmt.Println(2)
					fmt.Println(err)
					return
				}

				// fmt.Println(info.Data.Products[0].ID)
				for _, good := range info.Data.Products {
					fmt.Println(good.ID)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

}

func (a *Api) Getcarousels() {
	var out = make(map[int][]*Card)

	var wg sync.WaitGroup
	var curr int64

	var queue = goconcurrentqueue.NewFIFO()

	for i := 0; i < 1000; i++ {
		queue.Enqueue(i)
	}

	for i := 0; i < 100; i++ {
		// for j := 0; j < 12; j++ {
		wg.Add(1)
		go func() {
			for {
				c := context.Background()
				ctx, cancel := context.WithTimeout(c, time.Duration(time.Second*8))
				val, err := queue.Dequeue()
				if err != nil {
					cancel()
					break
				}

				ch := make(chan []*Card, 1)

				go a.check(ch, val.(int))

				select {
				case <-ctx.Done():
					fmt.Printf("Context cancelled: %v\n", ctx.Err())
				case result := <-ch:
					if len(result) > 0 {
						out[val.(int)] = result
						fmt.Println(val)
						for _, el := range result {
							fmt.Println(el.ID)
						}
					}
				}

				atomic.AddInt64(&curr, 1)
				go fmt.Println("Progress: " + strconv.FormatInt(curr, 10))
				cancel()
			}
			wg.Done()
		}()
		// }
	}
	wg.Wait()
	if err := PrettyPrint(out); err != nil {
		fmt.Println(err)
		return
	}
}

var ErrNotfound = errors.New("404")

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}

	f, e := os.Create("./parsed.json")
	if e != nil {
		fmt.Println(e)
		return
	}

	_, err = f.Write(b)
	return
}

func (a *Api) check(ch chan []*Card, id int) {
	cards, err := a.Search("@", id)
	if err != nil {
		fmt.Println(err)
		ch <- []*Card{}
		return
	}
	ch <- cards
}
