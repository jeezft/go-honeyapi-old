package utils

func UseAlphabet(a func(letter string) error) {
	for ch := 'a'; ch <= 'z'; ch++ {
		err := a(string(ch))
		if err != nil {
			continue
		}
	}
}
