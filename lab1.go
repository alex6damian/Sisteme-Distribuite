package main

import (
	"fmt"
	"time"
	)

// 3
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true
}

// 4
func primeste(chan1 chan<- string, msg string){
	// chan<- inseamna ca primeste doar
	chan1 <- msg
}

func trimite(chan1 <-chan string, chan2 chan<- string){
	// chan<- primeste
	// <-chan trimite
	// chan bidirectional
	msg := <- chan1
	chan2 <- msg
}

func main() {
	/*
	1. Să se implementeze o aplicație care sa transfere cu success un mesaj (ex.: salut) de la o rutină
	go (goroutine) la alta prin intermediul unui canal de comunicare.
	
	messages := make(chan string)

	go func() {
		messages <- "test"
	}()
	
	msg := <-messages
	fmt.Println(msg)
	*/

	/*
	2. Să se creeze o aplicație care primește două mesaje pe același canal de comunicare.

	messages := make(chan string, 2)

	messages <- "mesaj1"
	messages <- "mesaj2"

	fmt.Println(<-messages)
	fmt.Println(<-messages)
	*/

	/*
	3. Creați o aplicație care să notifice o altă rutină go că procesarea unei funcții a fost efectuată cu
	success.

	done := make(chan bool, 1)
	go worker(done)

	<-done
	*/

	/*
	4. Implementați o aplicație care să folosească un canal pentru primirea datelor și altul pentru
	trimiterea datelor.

	chan1 := make(chan string, 1)
	chan2 := make(chan string, 1)
	primeste(chan1, "Salut!")
	trimite(chan1, chan2)
	fmt.Println(<- chan2)
	*/

	/*
	5. Implementați o aplicație care să folosească instrucțiunea select cu scopul de a combina două
	rutine go (goroutines) pentru trimiterea de două sau mai multe mesaje.

	chan1 := make(chan string)
	chan2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		chan1 <- "primul mesaj"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		chan1 <- "al doilea mesaj"
	}()

	for range 2 {
		select {
		case mesaj1 := <- chan1:
			fmt.Println("Received:", mesaj1)
		case mesaj2 := <- chan2:
			fmt.Println("Received:", mesaj2)
		}
	}
	*/

	/*
	6. Să se implementeze o aplicație în care presupunem că executăm un apel extern care turnează
	rezultatul său pe canalul 1după 2 secunde.
	
	chan1 := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		chan1 <- "rezultat apel extern 1"
	}()

	select {
	case result := <- chan1:
		fmt.Println("Received:", result)
	case <- time.After(1 * time.Second):
		fmt.Println("timeout: apel extern prea lent")
	} // primeste timeout, rezultatul vine dupa 2 secunde iar timeout-ul este de 1 secunda

	chan2 := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		chan2 <- "rezultat apel extern 2"
	}()

	select {
	case result2 := <- chan2:
		fmt.Println("Received:", result2)
	case <- time.After(3 * time.Second):
		fmt.Println("timeout: apel extern prea lent")
	}
	*/
	
}