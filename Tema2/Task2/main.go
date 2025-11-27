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
	// rune for unicode support
	r := []rune(word)
	first, last := 0, len(r)-1
	for first < last {
		if r[first] != r[last] {
			return false
		}
		first++
		last--
	}
	return true
}

func main() {
	// input data
	input_data := [][]string{
		{"a1551a", "parc", "ana", "minim", "1pcl3"},
		{"calabalac", "tivit", "leu", "zece10", "ploaie", "9ana9"},
		{"lalalal", "tema", "papa", "ger"}}

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
	fmt.Printf("Medie palindroame: %.2f\n", average)
}
