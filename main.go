package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ascii-arcade/wish-template/app"
	"github.com/ascii-arcade/wish-template/language"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	version = "dev"
	host    = "localhost"
	port    = "23234"
)

func main() {
	langCode := cmp.Or(os.Getenv("ASCII_ARCADE_LANG"), "EN")
	lang := language.Languages[langCode]

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(app.TeaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error(lang.Get("ssh.could_not_start_server"), err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info(fmt.Sprintf(lang.Get("ssh.starting_server"), host, port))

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error(lang.Get("ssh.could_not_start_server"), err)
			done <- nil
		}
	}()

	<-done
	log.Info(lang.Get("ssh.stopping_server"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error(lang.Get("ssh.could_not_stop_server"), err)
	}
}
