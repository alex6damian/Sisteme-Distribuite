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
	"net/url"
	"os"
	"time"
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

// 3
func URLResourceManagement() {
	urlStr := "https://google.com:6969/path/to/resource?name=test&value=123"
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Scheme:", parsedURL.Scheme)
	fmt.Println("Host:", parsedURL.Host)
	fmt.Println("Path:", parsedURL.Path)

	// Query parameters
	queryParams := parsedURL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("Query Param: %s = %s\n", key, value)
		}
	}

	// getting a specific query parameter
	name := queryParams.Get("name")
	fmt.Println("Value of 'name' parameter:", name)
}

// 4
func timeFormats() {
	// default formats
	time := time.Now()
	fmt.Println("Default format ->", time)

	// custom formats
	// DD/MM/YYYY HH:MM:SS
	customFormat := time.Format("02/01/2006 15:04:05")
	fmt.Println("DD/MM/YYYY HH:MM:SS ->", customFormat)

	// YYYY-MM-DD
	customFormat = time.Format("2006-01-02")
	fmt.Println("YYYY-MM-DD ->", customFormat)

	// AM/PM format
	customFormat = time.Format("03:04:05 PM")
	fmt.Println("AM/PM format ->", customFormat)

	// month, day, year
	customFormat = time.Format("January 02, 2006")
	fmt.Println("Month, Day, Year ->", customFormat)
}

// 5
func getTimeUnits() {
	// methods to get seconds, milliseconds, nanoseconds
	// Unix(), UnixMilli(), UnixNano() extract time since epoch

	startTime := time.Now()
	time.Sleep(50 * time.Millisecond) // simulate some processing time
	endTime := time.Now()

	duration := endTime.Sub(startTime)

	fmt.Printf("Duration: %s\n\n", duration)

	// seconds returning float64, while milliseconds and nanoseconds return int64
	fmt.Printf("Duration in Seconds: %.4fs\n", duration.Seconds())
	fmt.Printf("Duration in Milliseconds: %dms\n", duration.Milliseconds())
	fmt.Printf("Duration in Nanoseconds: %dns\n", duration.Nanoseconds())
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

	fmt.Println("")
	URLResourceManagement()
	fmt.Println("")
	timeFormats()
	fmt.Println("")
	getTimeUnits()
}
