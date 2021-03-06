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

package cmd

import (
	"context"
	"github.com/mdhender/queenie/internal/config"
	"github.com/mdhender/queenie/internal/otohttp"
	"github.com/mdhender/queenie/internal/services/greeter"
	"github.com/mdhender/queenie/internal/services/solver"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var globalServe struct {
	debug struct {
		cors bool
	}
}

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "start the API server",
	Long:  `Start the API server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{ConfigFile: globalBase.ConfigFile}
		if err := config.Read(cfg); err != nil {
			log.Fatal(err)
		}
		if globalBase.VerboseFlag {
			log.Printf("[serve] %-30s == %q\n", "config", cfg.ConfigFile)
			log.Printf("[serve] %-30s == %q\n", "host", cfg.Server.Host)
			log.Printf("[serve] %-30s == %q\n", "port", cfg.Server.Port)
		}

		// start the server with the ability to shut it down gracefully
		// thanks to https://clavinjune.dev/en/blogs/golang-http-server-graceful-shutdown/.
		// TODO: this should be part of the server.Server implementation!
		log.Printf("server: todo: please move the shutdown logic to the server implementation!\n")

		// create a context that we can use to cancel the server
		ctx, cancel := context.WithCancel(context.Background())

		var options []otohttp.Option
		options = append(options, otohttp.WithContext(ctx))
		if globalServe.debug.cors {
			options = append(options, otohttp.WithDebugCors(true))
		}

		s, err := otohttp.NewServer(cfg, otohttp.Options(options...))
		if err != nil {
			log.Fatal(err)
		}

		greeter.RegisterGreeterService(s, greeter.Service{})
		if solverService, err := solver.NewService(); err != nil {
			log.Fatal(err)
		} else {
			solver.RegisterSolverService(s, solverService)
		}

		// run server in a go routine that we can cancel
		go func() {
			log.Printf("server: listening on %q\n", net.JoinHostPort(cfg.Server.Host, cfg.Server.Port))
			err := http.ListenAndServe(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), s)
			if err != http.ErrServerClosed {
				log.Fatalf("server: %v", err)
			}
		}()

		// catch signals to interrupt the server and shut it down
		chanSignal := make(chan os.Signal, 1)
		signal.Notify(chanSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-chanSignal
		log.Print("server: signal: interrupt: shutting down...\n")
		go func() {
			// in case the user is spraying us with interrupts...
			<-chanSignal
			log.Fatal("server: signal: kill: terminating...\n")
		}()

		// allow 5 seconds for a graceful shutdown
		ctxWithDelay, cancelNow := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelNow()
		// if that fails, panic!
		if err := s.Shutdown(ctxWithDelay); err != nil {
			panic(err)
		}
		log.Printf("server: stopped\n")

		// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
		cancel()

		defer os.Exit(0)
		return nil
	},
}

func init() {
	cmdServe.Flags().BoolVar(&globalServe.debug.cors, "debug-cors", false, "enable CORS debugging")

	cmdBase.AddCommand(cmdServe)
}
