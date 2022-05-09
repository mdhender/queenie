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

package definition

// SolverService lists the known words for a puzzle.
type SolverService interface {
	// Solve returns a solution.
	Solve(PuzzleRequest) SolutionResponse
}

// PuzzleRequest is the request object for SolverService.Solve
type PuzzleRequest struct {
	// Center letter is the required letter.
	// It must be a single, lower-case letter.
	// example: "c"
	Center string

	// Hex letters are the remaining six letters accepted in the solution.
	// It must be a string containing exactly six lower-case letters.
	// example: "hmnotu"
	Hex string
}

// SolutionResponse is the response object containing the list of known
// words that satisfy the puzzle.
type SolutionResponse struct {
	// Words is the list of known words that satisfy the puzzle.
	// example: ["cotton", "cottonmouth"]
	Words []string
}
