package kit

import (
	"errors"
	"fmt"
	"github.com/heptiolabs/healthcheck"
	_ "github.com/heptiolabs/healthcheck"
	"net/http"
	"time"
)

type HealthcheckConfig struct {
	Port string
}

type Healthcheck struct {
	hcHandler healthcheck.Handler
	port      string
	srv       *http.Server
}

func NewHealthCheck(cfg *HealthcheckConfig) *Healthcheck {
	return &Healthcheck{
		hcHandler: healthcheck.NewHandler(),
		port:      cfg.Port,
	}
}

func (h *Healthcheck) Start() {

	h.srv = &http.Server{
		Addr: fmt.Sprintf(":%s", h.port),
	}
	h.srv.Handler = h.hcHandler

	go func() {
	start:
		if err := h.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				time.Sleep(time.Second * 5)
				goto start
			}
			return
		}
	}()
}

func (h *Healthcheck) Stop() {
	if h.srv != nil {
		_ = h.srv.Close()
		h.srv = nil
	}
}
