package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func hangman(solution []string, wordarr []string) {
	imin := 0
	imax := 8
	fmt.Println("Good Luck, you have 10 attempts.")
	chance := 10
	for chance > 0 {
		equal := true
		for i := 0; i < len(solution); i++ {
			if solution[i] != wordarr[i] { // cas victoire
				equal = false
				break
			}
		}
		if equal == true {
			fmt.Println("Congrats !")
			return
		}
		for _, x := range solution {
			fmt.Print(string(x), " ")
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Choose: ")
		here := false
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		in := scanner.Text()
		for i, v := range wordarr {
			if in == v {
				solution[i] = string(v)
				here = true
			}
		}
		if here == false {
			chance--
			fmt.Println("Not present in the word, ", chance, " attempts remaining")
			if chance != 9 {
				imin += 8
				imax += 8
			}
		}
		dessin(chance, imin, imax)
	}
	fmt.Println("You lost, try again")
}

func dessin(chance int, imin int, imax int) {
	file, err := os.Open("hangman.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file) // LIRE UN MOT RDM
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())

	}
	for _, x := range text[imin:imax] {
		fmt.Println(string(x))
	}
}

func printn(word string, mi int, ma int) ([]string, []string) {
	solution := make([]string, len(word))
	wordarr := []string{}
	for _, v := range word {
		wordarr = append(wordarr, string(v))
	}
	for i := 0; i < len(word)/2-1; i++ { // une lettre random - n
		r := rand.Intn(ma-mi-1) + mi
		if solution[r] == "" {
			solution[r] = wordarr[r]
		}
	}
	for i := 0; i < len(solution); i++ {
		if solution[i] == "" {
			solution[i] = "_"
		}
	}
	return solution, wordarr
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Génère nb random
	min := 0
	max := 0
	file, err := os.Open("words.txt") //Ouvre .txt
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
		max++
	}
	random := rand.Intn(max-min) + min
	word := text[random]
	max = len(word)
	if err != nil {
		log.Fatalf("failed to open")
	}
	tmpsolution, tmpwordarr := (printn(word, min, max)) //Envoie et retrait de 2 tableaux [mot]-[solution]
	hangman(tmpsolution, tmpwordarr)                    //traitement
	file.Close()
}
