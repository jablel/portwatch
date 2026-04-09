package daemon

import (
	"context"
	"log"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

// Daemon runs the port monitoring loop.
type Daemon struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
}

// New creates a new Daemon with the provided configuration.
func New(cfg *config.Config, alerter *alert.Alerter) *Daemon {
	return &Daemon{
		cfg:     cfg,
		scanner: scanner.New(cfg),
		alerter: alerter,
	}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch starting — interval %s, ports %d-%d",
		d.cfg.Interval, d.cfg.PortRangeStart, d.cfg.PortRangeEnd)

	prev, err := state.Load(d.cfg.StateFile)
	if err != nil {
		log.Printf("no previous state found, starting fresh: %v", err)
		prev = &state.Snapshot{}
	}

	if err := d.tick(ctx, &prev); err != nil {
		return err
	}

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(ctx, &prev); err != nil {
				log.Printf("tick error: %v", err)
			}
		}
	}
}

func (d *Daemon) tick(ctx context.Context, prev *state.Snapshot) error {
	current, err := d.scanner.Scan(ctx)
	if err != nil {
		return err
	}

	snap := state.Snapshot{Ports: current}
	changes := state.Compare(*prev, snap)

	if err := d.alerter.Notify(changes); err != nil {
		log.Printf("alert error: %v", err)
	}

	if err := state.Save(d.cfg.StateFile, snap); err != nil {
		log.Printf("failed to save state: %v", err)
	}

	*prev = snap
	return nil
}
