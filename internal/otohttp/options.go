/*
 * queenie - a spelling bee helper
 * Copyright (C) 2022 Michael D Henderson
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
 * Much of the Oto HTTP source is pulled from Pace Software's oto project
 * (see https://github.com/pacedotdev/oto/tree/main/otohttp). That code
 * (and my changes to their code) are released under the following MIT License:
 *
 *  Copyright (c) 2021 Pace Software Ltd
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 *
 */

package otohttp

import (
	"context"
	"net/http"
)

// Option allows us to pass in options when creating a new server.
type Option func(*Server) error

// Options turns a list of Option instances into an Option.
func Options(opts ...Option) Option {
	return func(s *Server) error {
		for _, opt := range opts {
			if err := opt(s); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithBasepath changes the default from `/oto/`.
func WithBasepath(path string) Option {
	return func(s *Server) (err error) {
		s.Basepath = path
		return nil
	}
}

// WithContext adds a context to the server so that we can shut it down gracefully.
func WithContext(ctx context.Context) Option {
	return func(s *Server) (err error) {
		s.ctx = ctx
		return nil
	}
}

// WithNotFound changes the default not found handler.
func WithNotFound(h http.Handler) Option {
	return func(s *Server) (err error) {
		s.NotFound = h
		return nil
	}
}

// WithOnErr changes the default error handler.
func WithOnErr(fn func(http.ResponseWriter, *http.Request, error)) Option {
	return func(s *Server) (err error) {
		s.OnErr = fn
		return nil
	}
}
