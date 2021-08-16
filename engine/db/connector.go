package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type Connector struct {
	config *Config
	driver string
	db     *sqlx.DB
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) init() {
	c.config = NewConfig()
	c.driver = os.Getenv("DB_DRIVER")
}

func (c *Connector) SetConnectionX(db *sqlx.DB) {
	c.db = db
	c.driver = db.DriverName()
}

func (c *Connector) SetConnection(db *sql.DB, driver string) {
	c.db = sqlx.NewDb(db, driver)
}

func (c *Connector) GetDriver() string {
	return c.driver
}

func (c *Connector) GetConnection() *sqlx.DB {
	return c.db
}

func (c *Connector) Connect() error {
	var err error
	// If connection was already set just return the instance
	if c.db != nil {
		return err
	}

	c.init() // initialize a config for the new DB connection
	connectionString, ok := c.getConnectionString()
	if !ok {
		return fmt.Errorf("driver '%s' is not supported", c.driver)
	}
	c.db, err = sqlx.Connect(c.driver, connectionString)
	return err
}

func (c *Connector) Close() {
	if c.db != nil {
		_ = c.db.Close()
	}
}

func (c *Connector) getConnectionString() (string, bool) {
	switch c.driver {
	case "mysql":
		return c.getMysqlConnectionString(), true
	case "postgres":
		return c.getPostgresConnectionString(), true
	default:
		return "", false
	}
}

func (c *Connector) getMysqlConnectionString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		c.config.user, c.config.password, c.config.host, c.config.port, c.config.name)
}

func (c *Connector) getPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.config.host, c.config.port, c.config.user, c.config.password, c.config.name)
}
