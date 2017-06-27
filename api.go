package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

// Dictionary represents a list of words
type Dictionary struct {
	elements []string
}

// Select will select a word from a dictionary
func (d *Dictionary) Select(i int) string {
	return d.elements[i]
}

// SelectRandom will select a random word from the dictionary
func (d *Dictionary) SelectRandom() string {
	index := rand.Intn(len(d.elements))
	fmt.Printf("%d\n", index)
	return d.Select(index)
}

// Load will load a new Dictionary
func Load(fileName string) (*Dictionary, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	s := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, string(data))
	return &Dictionary{cleanse(strings.Split(s, " "))}, nil
}

func cleanse(s []string) []string {
	news := []string{}
	for _, x := range s {
		if strings.Replace(x, " ", "", -1) == "" {
			continue
		}
		news = append(news, x)
	}
	return news
}

// Run runs the API instance
func main() {
	rand.Seed(time.Now().UnixNano())

	adjDict, err := Load("api/adj_dict.txt")

	if err != nil {
		log.Fatal(err)
	}

	animalDict, err := Load("api/animal_dict.txt")

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, generateName(adjDict, animalDict))
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

func generateName(dict ...*Dictionary) string {
	s := ""
	for i, d := range dict {
		s += d.SelectRandom()

		if i != len(dict)-1 {
			s += "-"
		}
	}
	return s
}
