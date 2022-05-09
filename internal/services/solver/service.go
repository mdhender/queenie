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

package solver

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Service struct {
	dict    map[string]bool
	invalid map[string]bool
	valid   map[string]bool
	checks  map[string]bool
	words   []string
}

func NewService() (Service, error) {
	s := Service{}

	// load the words file sourced from https://github.com/dwyl/english-words and other places
	var err error
	if s.dict, err = loadWords("wordlist.txt"); err != nil {
		return Service{}, err
	} else if s.invalid, err = loadWords("invalid.txt"); err != nil {
		return Service{}, err
	} else if s.valid, err = loadWords("valid.txt"); err != nil {
		return Service{}, err
	} else if s.checks, err = loadWords("checks.txt"); err != nil {
		return Service{}, err
	}

	for word := range s.dict {
		s.words = append(s.words, word)
	}
	sort.Strings(s.words)

	return s, nil
}

func (s Service) Solve(ctx context.Context, request PuzzleRequest) (*SolutionResponse, error) {
	// consider rejecting unknown fields (https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)
	if request.Center == "" {
		return nil, errors.New("missing 'center'")
	} else if request.Hex == "" {
		return nil, errors.New("missing 'hex'")
	}

	var centerLetter rune
	for i, r := range request.Center {
		if i > 0 || !unicode.IsLetter(r) {
			return nil, errors.New("invalid 'center'")
		}
		centerLetter = unicode.ToLower(r)
	}

	var hexLetters []rune
	for i, r := range request.Hex {
		if i > 6 || !unicode.IsLetter(r) {
			return nil, errors.New("invalid 'hex'")
		}
		r = unicode.ToLower(r)
		if r == centerLetter {
			return nil, errors.New("duplicate 'hex'")
		}
		for _, h := range hexLetters {
			if r == h {
				return nil, errors.New("duplicate 'hex'")
			}
		}
		hexLetters = append(hexLetters, r)
	}
	if len(hexLetters) != 6 {
		return nil, errors.New("invalid 'hex'")
	}

	// the first letter in it must be the puzzle's center letter.
	letters := [7]rune{centerLetter, hexLetters[0], hexLetters[1], hexLetters[2], hexLetters[3], hexLetters[4], hexLetters[5]}

	var words []string
	for _, word := range s.words {
		if !strings.ContainsRune(word, centerLetter) {
			continue
		}
		w := word
		// remove all the allowed characters from the word
		for _, ch := range letters {
			w = strings.ReplaceAll(w, string(ch), "")
		}
		// result must be empty (ie, must not contain any other letters)
		if len(w) == 0 {
			words = append(words, word)
		}
	}
	//sort.Strings(words)

	return &SolutionResponse{
		Words: words,
	}, nil
}

func loadWords(filename string) (map[string]bool, error) {
	raw, err := os.ReadFile(filename)
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
