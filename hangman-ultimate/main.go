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
	asciiletter(textascii, hgn)
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
			hgn.Word = hgn.Solution
			asciiletter(textascii, hgn)
			os.Remove("resultat.txt")
			return
		}
		/*for _, x := range hgn.word { //print mot incomplet
			fmt.Print(string(x), " ")
		}*/
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Choose: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text() //prend en entrée un mot ou une lettre
		here := false

		if stopsave(input, hgn) == "STOP" {
			break
		}

		if len(input) > 1 && input != "STOP" { //Traite input == mot
			if input == strsolution {
				fmt.Println("Congrats !")
				hgn.Word = hgn.Solution
				asciiletter(textascii, hgn)
				os.Remove("resultat.txt")
				return
			} else {
				hgn.Chance -= 2
				if hgn.Chance == 8 {
					hgn.Drawingpos += 8
				} else {
					hgn.Drawingpos += 16
				}
				fmt.Println("Not present in the word, ", hgn.Chance, " attempts remaining")
			}
		} else { //Traite input == lettre
			for _, v := range inputletter {
				if v == input {
					fmt.Println("Error letter already submitted")
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
		}

		if here == false && len(input) <= 1 { //pas présent print chance-1 et dessin correspondant
			hgn.Chance--
			fmt.Println("Not present in the word, ", hgn.Chance, " attempts remaining")
			if hgn.Chance != 9 {
				hgn.Drawingpos += 8
			}
		}
		asciiletter(textascii, hgn)
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

//Prend le mot en parametre et renvoie le résultat en lettres acscii passant en paramètres
func asciiletter(ascii []string, hgn hangm) {
	file, err := os.OpenFile("resultat.txt", os.O_CREATE|os.O_WRONLY, 0600)
	defer file.Close() // on ferme automatiquement à la fin de notre programme
	//L'objectif ici est de donner une valeur (un dessin) à la letrre "OS.arg[:1]"
	slicealpha := []string{}
	alpha := "abcdefghijklmnopqrstuvwxyz"
	count := 0
	essaie := []string{}
	//Ici min et max correspondent aux index séparant les lettre ascii que nous imprimerons par la suite
	for _, x := range alpha {
		slicealpha = append(slicealpha, string(x))
	}
	for i := range hgn.Word {
		/*
			Si la lettre recu correspond à une dès lettre de lalphabet contenue dans le slice (slicealpha)
			alors (count) calcul la position de la lettre en question
		*/
		for j := range slicealpha {
			count++
			if hgn.Word[i] == slicealpha[j] {
				min := (count * 9) + 577
				max := (count * 9) + 586

				for _, x := range ascii[min:max] {
					_, err = file.WriteString("\n")
					_, err = file.WriteString(string(x))
					// écrire dans le fichier "resultat.txt" la lettre correspondante
					if err != nil {
						panic(err)
					}
				}
			}
		}
		if hgn.Word[i] != hgn.Solution[i] {
			for _, x := range ascii[116:125] {
				_, err = file.WriteString("\n")
				_, err = file.WriteString(string(x))
				// écrire dans le fichier "resultat.txt" le caractère "_"
				if err != nil {
					panic(err)
				}
			}

		}
		count = 0
	}
	counter := 9
	counter2 := 1
	for j := 0; j < 9; j++ { //Recupère l'index+1 \n de chacune des lettres dans "resultat.txt"
		for i := len(hgn.Solution); i > 0; i-- {
			counter2++
			n := 0
			n = 9*counter2 - counter
			filer, _ := os.Open("resultat.txt")
			scanner := bufio.NewScanner(filer)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				essaie = append(essaie, scanner.Text())
			}
			fmt.Print(essaie[n])
		}
		if j != 9 {
			fmt.Print("\n")
		}
		counter2 = 0
		counter--
	}
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
