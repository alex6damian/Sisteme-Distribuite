package main

import (
	"fmt"
	"strings"
	"sync"
)

// pair for mapping
type pair struct {
	key int
	val int
}

func worker(i int, words []string, wg *sync.WaitGroup, output chan<- pair) {
	// close waitgroup when done
	defer wg.Done()

	// iterate over words
	for _, word := range words {
		if valid(word) {
			// send pair to output channel
			output <- pair{key: i, val: 1}
		}
	}
}

func valid(word string) bool {
	word = strings.ToUpper(word)
	vow, cons := 0, 0
	for _, chr := range word {
		if chr == 'A' || chr == 'E' || chr == 'I' || chr == 'O' || chr == 'U' {
			vow++
		} else {
			cons++
		}
	}
	return vow%2 == 0 && cons%3 == 0
}

func main() {

	// input data
	input_data := [][]string{
		{"aabbb", "ebep", "blablablaa", "hijk", "wsww"},
		{"abba", "eeeppp", "cocor", "ppppppaa", "qwerty", "acasq"},
		{"lalala", "lalal", "papapa", "papap"}}

	// waitgroup for goroutines
	var wg sync.WaitGroup

	// channel for output pairs
	output := make(chan pair, 50)

	for i, arr := range input_data {
		// increment waitgroup counter and start goroutine for each array
		wg.Add(1)
		go worker(i, arr, &wg, output)
	}

	go func() {
		// wait for all goroutines to finish
		wg.Wait()
		// close output channel
		close(output)
	}()

	// shuffle and reduce
	cnt := make(map[int]int)
	for pair := range output {
		cnt[pair.key] += pair.val
	}

	// calculate results
	groupsNumber := len(input_data)
	sum := 0
	for i := range cnt {
		sum += cnt[i]
	}
	average := float64(sum) / float64(groupsNumber)

	fmt.Printf("Medie cuvinte cu nr par de vocale si nr divizibil cu 3 de consoane: %.2f\n", average)
}
