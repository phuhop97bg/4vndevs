package user_storage

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

// MySQLConfig type;
type MySQLConfig struct {
	Name        string // for trace
	Environment string
	DSN         string // data source name
	Active      int    // pool
	Idle        int    // pool
	Lifetime    time.Duration
}

// NewSqlxDB type;
func NewSqlxDB(c *MySQLConfig) (db *sqlx.DB) {
	db, err := sqlx.Connect("mysql", c.DSN)
	if err != nil {
		glog.V(1).Infof("Connect db error: %s", err)
	}

	glog.V(3).Infof("NewSqlxDB: %+v", db.Stats())
	db.SetConnMaxLifetime(c.Lifetime * time.Second)
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	return
}
