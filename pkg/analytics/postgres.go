package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Event struct {
	Timestamp time.Time
	Endpoint  string
	ClientIP  string
	Status    int
	LatencyMs float64
	Country   *string
	Continent *string
}

type Collector struct {
	db     *sql.DB
	geoip  *GeoIP
	events chan Event
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewCollector(postgresURL, geoipPath string) (*Collector, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	var geoip *GeoIP
	if geoipPath != "" {
		geoip, err = NewGeoIP(geoipPath)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("open geoip: %w", err)
		}
	}

	c := &Collector{
		db:     db,
		geoip:  geoip,
		events: make(chan Event, 4096),
		done:   make(chan struct{}),
	}
	c.wg.Add(1)
	go c.worker()
	return c, nil
}

func (c *Collector) DB() *sql.DB {
	return c.db
}

func (c *Collector) Record(e Event) {
	if c.geoip != nil {
		e.Country, e.Continent = c.geoip.Lookup(e.ClientIP)
	}
	select {
	case c.events <- e:
	default:
		// channel full, drop event
	}
}

func (c *Collector) Close() {
	close(c.done)
	c.wg.Wait()
	if c.geoip != nil {
		c.geoip.Close()
	}
	c.db.Close()
}

func (c *Collector) worker() {
	defer c.wg.Done()

	batch := make([]Event, 0, 100)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case e, ok := <-c.events:
			if !ok {
				// channel closed
				if len(batch) > 0 {
					c.flush(batch)
				}
				return
			}
			batch = append(batch, e)
			if len(batch) >= 100 {
				c.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.flush(batch)
				batch = batch[:0]
			}
		case <-c.done:
			// drain remaining events
			for {
				select {
				case e := <-c.events:
					batch = append(batch, e)
				default:
					if len(batch) > 0 {
						c.flush(batch)
					}
					return
				}
			}
		}
	}
}

func (c *Collector) flush(batch []Event) {
	if len(batch) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("[ANALYTICS]: Failed to begin tx: %v\n", err)
		return
	}

	var b strings.Builder
	b.WriteString("INSERT INTO events (timestamp, endpoint, client_ip, status, latency_ms, country, continent) VALUES ")

	args := make([]interface{}, 0, len(batch)*7)
	for i, e := range batch {
		if i > 0 {
			b.WriteString(",")
		}
		offset := i * 7
		fmt.Fprintf(&b, "($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7)
		args = append(args, e.Timestamp, e.Endpoint, e.ClientIP, e.Status, e.LatencyMs, e.Country, e.Continent)
	}

	if _, err := tx.ExecContext(ctx, b.String(), args...); err != nil {
		tx.Rollback()
		fmt.Printf("[ANALYTICS]: Failed to insert batch: %v\n", err)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("[ANALYTICS]: Failed to commit: %v\n", err)
	}
}

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS events (
			id         BIGSERIAL PRIMARY KEY,
			timestamp  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			endpoint   TEXT NOT NULL,
			client_ip  TEXT NOT NULL,
			status     INTEGER NOT NULL,
			latency_ms DOUBLE PRECISION,
			country    TEXT,
			continent  TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events (timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_events_endpoint_timestamp ON events (endpoint, timestamp)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
