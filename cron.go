package main

import (
	"fmt"

	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
)

type cron struct {
	SessionExpirationDays int64
	Database              *sqlx.DB
}

func (c *cron) startJobs() {
	go func() { <-gocron.Start() }()
}

func (c *cron) pruneExpiredSessions() {
	// Actually hard-delete "session users"
	// rather than retain data, to satisfy GDPR
	c.Database.MustExec(fmt.Sprintf(`
		DELETE FROM user WHERE 1=1
			AND email NOT LIKE "%%@\%%"
			AND created_at < DATE_SUB(CURDATE(), INTERVAL %d DAY)
	`, c.SessionExpirationDays))
}
