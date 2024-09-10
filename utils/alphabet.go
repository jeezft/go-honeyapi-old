package utils

import "sync"

func UseAlphabet(a func(letter string) error) {
	var wg sync.WaitGroup
	for ch := 'a'; ch <= 'z'; ch++ {
		wg.Add(1)
		go func(ch rune) {
			err := a(string(ch))
			wg.Done()
			if err != nil {
				return
			}
		}(ch)

	}
	wg.Add(2)

	go func() {
		_ = a("%D0%B0-%D1%8F")
		wg.Done()
	}()

	go func() {
		_ = a("123")
		wg.Done()
	}()

	wg.Wait()
}
