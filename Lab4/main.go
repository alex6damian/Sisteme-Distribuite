package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// opening the file and reading data using Unmarshal into a map
func getData(filename string) (map[string]float64, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var fructe map[string]float64
	err = json.Unmarshal(file, &fructe)
	if err != nil {
		return nil, err
	}
	return fructe, nil
}

func main() {

	// getting the data and printing it
	fructe, err := getData("cantitati_fructe.json")
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	for fruct, cantitate := range fructe {
		fmt.Printf("Fruct: %s, Cantitate: %.2f kg\n", fruct, cantitate)
	}

}
