/*
Cerinte:
I.
    1. Să se creeze o aplicație care să verifice dacă un fișier există sau nu într-o anumită locație.
    2. Să se îmbunătățească programul creat anterior prin adăugarea de excepții și erori, astfel
    încât dacă fișierul există, programul să poată continua execuția, altfel să afișeze o eroare
    și să oprească continuarea programului.
    3. Scieți o aplicație care să citească un fișier și conținutul fișierului să fie salvat într-o
    variabilă de tip string.
    4. Stocați conținutul unui fișier într-un vector, prin citirea fiecărei linii în parte și salvarea
    acesteia în vector.
    5. Scrieți un program care să creeze un fișier în care să se salveze anumite informații.
    6. Scrieți un program care să ajute la procesul de redenumire a unui fișier pe hard disk.
II.
    1. Scrieți o aplicație pentru validarea expresiilor/a datelor de intrare.
III.
    1. Scrieți o aplicație care să codeze/decodeze în format JSON diferite date și să afișeze datele în
    format JSON.    
IV.
    1. Implementați o aplicație de tip HTTP Client care să returneze rezultatul răspunsului unei
    interogări tip HTTP pentru un anumit site (e.g., http://www.google.com)
    2. Implementați un server HTTP folosind pachetul net/http și care să permită apeluri pentru
    rute predefinite (e.g., /hello)
*/
package main

import (
    "fmt"
    "errors"
    "os"
    "log"
    "bufio"
    "regexp"
)

// I.1/2
func fileExists(path string) bool {
    _, err := os.Stat(path)
    // os.Stat returneaza informatii despre fisier sau o eroare daca fisierul nu exista
    if errors.Is(err, os.ErrNotExist) {
        log.Fatalf("Eroare: Fisierul %s nu exista.\n", path)
    }
    return true
}
// I.3
func readFromFile(path string) string {
    if fileExists(path) {
        dat, _ := os.ReadFile("exemplu.txt")
        return string(dat)
    }
    return ""
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}
// I.4
func readMultipleLines(path string) []string {
    var lines []string
    if fileExists(path) {
        f, err := os.Open(path)
        check(err)
        defer f.Close() // inchide fisierul la finalul functiei

        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            lines = append(lines, scanner.Text())
        }
        if err := scanner.Err(); err != nil { // verific erorile de citire
            log.Fatal(err)
        }
    }
    return lines
}
// I.5
func createFile(filename string, content string) {
    f, err := os.Create(filename)
    check(err)
    defer f.Close()
    f.WriteString(content)
}

// I.6
func renameFile(oldName string, newName string) {
    if !fileExists(oldName) {
        fmt.Printf("Eroare: Fisierul %s nu exista.\n", oldName)
        return
    }
    err := os.Rename(oldName, newName)
    check(err)
    fmt.Printf("Fisierul %s a fost redenumit in %s.\n", oldName, newName)
}

// II.1
func validatePhoneNumber(phone string) bool {
    pattern := `^07[0-9]{8}$`
    regex := regexp.MustCompile(pattern)
    return regex.MatchString(phone)
}

func main() {
    path := "lab2.go"
    if fileExists(path) {
        fmt.Printf("Fisierul %s exista.\n", path)
    }

    content := readFromFile("exemplu.txt")
    fmt.Println("Continutul fisierului:", content)

    lines := readMultipleLines("exemplu.txt")
    fmt.Println("Liniile din fisier:", lines)

    createFile("fisier_generat.txt", "Acesta este un fisier creat de program.\n")

    renameFile("fisier_generat.txt", "fisier_redenumit.txt")

    // accepta doar nr de 10 cifre care incep cu 07
    numere := []string{"0712345678", "1234567890", "0798765432"}
    for _, numar := range numere {
        if validatePhoneNumber(numar) {
            fmt.Printf("Numarul %s este valid.\n", numar)
        } else {
            fmt.Printf("Numarul %s nu este valid.\n", numar)
        }
    }
}
