package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type hangm struct {
	Solution   []string //mot solution
	Word       []string //mot incomplet
	Chance     int      //nb chance
	Drawingpos int      //dessin correspondant au nb de chance(s)
}

/*Prend la lettre ou le mot choisi et vérifie si le mot est complet sinon renvoie le dessin
correspondant etle nombre de chance restant */
func hangman(hgn hangm) {
	textascii := []string{}
	file, err := os.Open("standard.txt")
	if err != nil {
		log.Fatalf("failed to open standard.txt")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		textascii = append(textascii, scanner.Text()) // Slice ascii
	}
	strsolution := ""
	for _, v := range hgn.Solution { //Incrémentation de la solutin
		strsolution += string(v)
	}
	inputletter := []string{}
	if hgn.Chance == 10 {
		fmt.Println("Good Luck, you have", hgn.Chance, "attempts.")
	} else {
		fmt.Println("Welcome Back, you have", hgn.Chance, "attempts remaining.")
	}
	for hgn.Chance > 0 {
		equal := true
		for i := 0; i < len(hgn.Word); i++ { //vérifie si le mot est complet
			if hgn.Word[i] != hgn.Solution[i] {
				equal = false
				break
			}
		}
		if equal == true { //si le mot est complet victoire et arrêt
			fmt.Println("Congrats !")
			os.Remove("resultat.txt")
			return
		}
		for _, x := range hgn.Word { //print mot incomplet
			fmt.Print(string(x), " ")
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Choose: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text() //prend en entrée une lettre
		here := false

		if stopsave(input, hgn) == "STOP" {
			break
		}
		for _, v := range inputletter {
			if v == input {
				fmt.Println("Error letter already submitted")
				here = true
			} else {
				for i, v := range hgn.Solution {
					if input == v {
						hgn.Word[i] = string(v)
						here = true
					}
				}
			}
		}
		inputletter = append(inputletter, input)

		if here == false && len(input) <= 1 { //pas présent print chance-1 et dessin correspondant
			hgn.Chance--
			fmt.Println("Not present in the word, ", hgn.Chance, " attempts remaining")
			if hgn.Chance != 9 {
				hgn.Drawingpos += 8
			}
		}
		dessin(hgn)
	}
	if hgn.Chance == 0 {
		fmt.Println("You lost, try again")
	}
	os.Remove("resultat.txt")
}

func stopsave(input string, hgn hangm) string {
	if input == "STOP" { //START AND STOP
		save, _ := json.Marshal(hgn)
		f, _ := os.Create("save.json")
		f.Write(save)
		fmt.Println("Game Saved in save.json.")
		return "STOP"
	}
	return ""
}

//print dessin correspondant
func dessin(hgn hangm) {
	var text []string
	file, err := os.Open("hangman.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file) //lis le fichier qui contient les dessins
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	for _, x := range text[hgn.Drawingpos : hgn.Drawingpos+8] { //Print le dessin ,via la dernière position drawingpos(l'index) défini dans la struct +8
		fmt.Println(string(x))
	}
}

//initialisation du mot incomplet et de sa solution
func printn(word string, mi int, ma int) ([]string, []string) {
	wordarr := make([]string, len(word))
	solution := []string{}
	for _, v := range word { //met le mot complet dans solution
		solution = append(solution, string(v))
	}
	for i := 0; i < len(word)/2-1; i++ { //Révèle n lettres random du mot où n est len(word) / 2 - 1
		r := rand.Intn(ma-mi-1) + mi
		if wordarr[r] == "" {
			wordarr[r] = solution[r]
		}
	}
	for i := 0; i < len(wordarr); i++ { //Remplace chaque case vide par un "_"
		if wordarr[i] == "" {
			wordarr[i] = "_"
		}
	}
	return wordarr, solution
}

func main() {
	//_______________________________________
	aff, _ := os.Open("affichage.txt")
	scan := bufio.NewScanner(aff)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		fmt.Println(scan.Text())
	}
	fmt.Println("")
	//_______________________________________
	args := os.Args[1:]
	var text []string
	min := 0
	max := 0
	rand.Seed(time.Now().UnixNano()) // Génère nb random
	file, err := os.Open("words.txt")
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text = append(text, scanner.Text())
		max++ //borne max random
	}
	random := rand.Intn(max-min) + min
	word := text[random]
	max = len(word)
	if err != nil {
		log.Fatalf("failed to open")
	}
	tmpword, tmpsolution := (printn(word, min, max))
	hgn := hangm{tmpsolution, tmpword, 10, 0}

	if len(args) > 1 { //Récupération de la sauvegarde si flag trouvé
		if args[0] == "--startWith" && args[1] == "save.txt" {
			save, _ := ioutil.ReadFile("save.json")
			err = json.Unmarshal([]byte(save), &hgn)
			if err != nil {
				log.Fatalf("failed to encode game")
			}
		}
	}
	hangman(hgn)
}
