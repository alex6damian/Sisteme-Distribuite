/*
1. Scrieți o aplicație care să codeze/decodeze în format JSON diferite date și să afișeze datele în
format JSON.
2. Implementați o aplicație care să genereze un document (fișier) XML pentru gestionarea
cafelei și originii acesteia.
3. Scrieți o aplicație care să gestioneze resursele și parametrii unui URL.
4. Scrieți o aplicație care să permită gestionarea și afișarea diferitelor formate de dată și timp.
5. Implementați mecanisme Go care să ajute la obținerea de secunde, milisecunde și
nanosecunde.
*/
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile(filename string) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
}

type Coffee struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Plant struct {
	XMLName xml.Name `xml:"plant"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:"name"`
	Origin  string   `xml:"origin"`
}

func (p Plant) String() string {
	return fmt.Sprintf("Plant[Id=%d, Name=%s, Origin=%s]", p.Id, p.Name, p.Origin)
}

func encodeDecodeJSON() {
	// Encoding
	obiect := &Coffee{
		Name:  "Espresso",
		Price: 4.5}
	result, _ := json.Marshal(obiect)
	fmt.Println("Encoded JSON:", string(result))

	intB, _ := json.Marshal(123)
	fmt.Println("Encoded JSON Integer:", string(intB))

	// Decoding
	var dat map[string]interface{}
	if err := json.Unmarshal(result, &dat); err != nil {
		panic(err)
	}
	fmt.Println("Decoded JSON:", dat)

	var num int
	if err := json.Unmarshal(intB, &num); err != nil {
		panic(err)
	}
	fmt.Println("Decoded JSON Integer:", num)
}

func main() {
	encodeDecodeJSON()

	createFile("cafea.xml")

	plant := &Plant{Id: 1, Name: "Arabica", Origin: "Ethiopia"}
	plant2 := &Plant{Id: 2, Name: "Robusta", Origin: "Vietnam"}
	plant3 := &Plant{Id: 3, Name: "Liberica", Origin: "Philippines"}
	plants := []Plant{*plant, *plant2, *plant3}
	out, _ := xml.MarshalIndent(plants, " ", " ")
	xmlContent := xml.Header + string(out)
	err := os.WriteFile("cafea.xml", []byte(xmlContent), 0644)
	if err != nil {
		panic(err)
	}
}
