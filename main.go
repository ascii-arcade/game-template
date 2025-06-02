package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ascii-arcade/wish-template/app"
	"github.com/ascii-arcade/wish-template/config"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(config.Host, config.Port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(app.TeaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error(config.Language.Get("ssh.could_not_start_server"), err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info(fmt.Sprintf(config.Language.Get("ssh.starting_server"), config.Host, config.Port))

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error(config.Language.Get("ssh.could_not_start_server"), err)
			done <- nil
		}
	}()

	<-done
	log.Info(config.Language.Get("ssh.stopping_server"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error(config.Language.Get("ssh.could_not_stop_server"), err)
	}
}
