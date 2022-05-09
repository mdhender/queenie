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
	"unicode"
)

type Service struct{}

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

	return &SolutionResponse{
		Words: []string{"cotton", "cottonmouth", string(centerLetter)},
	}, nil
}
