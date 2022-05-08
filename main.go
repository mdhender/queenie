/*
 * queenie - a spelling bee helper
 * Copyright (C) 2022 Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

// Package main implements the entry point for the Queenie server.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func main() {
	started := time.Now()

	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	if err := run(); err != nil {
		log.Println(err)
	}

	elapsed := time.Now().Sub(started)
	log.Printf("elapsed time: %+v\n", elapsed)
}

func run() (err error) {
	s := &server{}

	// load the words file sourced from https://github.com/dwyl/english-words and other places
	if s.dict, err = loadwords("wordlist.txt"); err != nil {
		return err
	} else if s.invalid, err = loadwords("invalid.txt"); err != nil {
		return err
	} else if s.valid, err = loadwords("valid.txt"); err != nil {
		return err
	} else if s.checks, err = loadwords("checks.txt"); err != nil {
		return err
	}

	// add valid to dictionary
	for word := range s.valid {
		s.dict[word] = true
	}

	router := http.NewServeMux()
	router.Handle("/", s.index())

	port := 8080
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.ListenAndServe()
}

type server struct {
	dict    map[string]bool
	invalid map[string]bool
	valid   map[string]bool
	checks  map[string]bool
}

func (s *server) index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path[1:]
		if len(route) != 7 {
			_, _ = fmt.Fprintf(w, "invalid route")
			return
		}

		_, _ = fmt.Fprintln(w, "<html>")
		_, _ = fmt.Fprintln(w, "<body>")

		// letters should be a command line flag.
		// the first letter in it must be the puzzle's required letter.
		letters := []byte(route)

		requiredLetter := rune(letters[0])

		matches := make(map[string]bool)
		for word := range s.dict {
			if !strings.ContainsRune(word, requiredLetter) {
				continue
			}
			w := word
			// remove all the allowed characters from the word
			for _, ch := range letters {
				w = strings.ReplaceAll(w, string(ch), "")
			}
			// result must be empty (ie, must not contain any other letters)
			if len(w) != 0 {
				continue
			}
			matches[word] = true
		}

		_, _ = fmt.Fprintln(w, "<ul>")
		if len(s.checks) != 0 {
			for word := range matches {
				if s.checks[word] {
					continue
				}
				_, _ = fmt.Fprintf(w, "<li>%s</li>", word)
			}
		} else {
			for word := range matches {
				if s.invalid[word] {
					continue
				}
				if s.dict[word] {
					_, _ = fmt.Fprintf(w, "<li> %s</li>", word)
					continue
				}
				if s.valid[word] {
					_, _ = fmt.Fprintf(w, "<li>+%s</li>", word)
					continue
				}
			}
		}
		_, _ = fmt.Fprintln(w, "</ul>")
		_, _ = fmt.Fprintln(w, "</body>")
		_, _ = fmt.Fprintln(w, "</html>")
	})
}

func loadwords(filename string) (map[string]bool, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	words := make(map[string]bool)
	for _, word := range strings.Split(string(raw), "\n") {
		// must be at least four characters
		if len(word) < 4 {
			continue
		}
		words[strings.ToLower(word)] = true
	}

	return words, nil
}

func printmap(words []string, f func(word string)) {
	sort.Strings(words)
	for _, word := range words {
		f(word)
	}
}
