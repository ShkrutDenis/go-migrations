package migrator

import (
	"time"
)

type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Batch     int       `db:"batch"`
	CreatedAt time.Time `db:"created_at"`
}
