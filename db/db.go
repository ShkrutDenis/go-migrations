package db

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Connector struct {
	config     *ConnectionConfig
	driver     string
	connection *sqlx.DB
}

type ConnectionConfig struct {
	user       string
	password   string
	connection string
	host       string
	port       string
	name       string
}

func NewConnector() *Connector {
	config := NewConfig()
	return &Connector{
		config: config,
		driver: os.Getenv("DB_DRIVER"),
	}
}

func NewConfig() *ConnectionConfig {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if host == "" && port == "" {
		connectionParts := strings.Split(os.Getenv("DB_CONNECTION"), ":")
		if len(connectionParts) == 2 {
			host, port = connectionParts[0], connectionParts[1]
		}
	}

	return &ConnectionConfig{
		user:       os.Getenv("DB_USER"),
		password:   os.Getenv("DB_PASSWORD"),
		connection: os.Getenv("DB_CONNECTION"),
		host:       host,
		port:       port,
		name:       os.Getenv("DB_NAME"),
	}
}

func (c *Connector) Connect() (*sqlx.DB, error) {
	connectionString, ok := c.getConnectionString()
	if !ok {
		return nil, errors.New("driver was not provided")
	}
	return sqlx.Connect(os.Getenv("DB_DRIVER"), connectionString)
}

func (c *Connector) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return errors.New("connection not exist")
}

func (c *Connector) ChangeDB(dbName string) {
	c.config.name = dbName
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
	return fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true",
		c.config.user, c.config.password, c.config.connection, c.config.name)
}

func (c *Connector) getPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.config.host, c.config.port, c.config.user, c.config.password, c.config.name)
}
