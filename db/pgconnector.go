package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// DB Logger for queries
type dbLogger struct{}

// Dummy Function to execute before a query
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// Log the generated query after execution
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}

//Adiitional constants for pg connection
const (
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
	PoolSize     = 10
	MinIdleConns = 10
	AppName      = "neo-app"
)

// DB Model for a NEO
type Neo struct {
	Id                             string    `pg:",pk"`
	Name                           string    `pg:"name"`
	NasaJplUrl                     string    `pg:"jpl_url"`
	IsPotentiallyHazardousAsteroid bool      `pg:"hazardous,notnull,use_zero"`
	Date                           time.Time `pg:"date"`
}

// Type with config options for the db connection
type DbConfig struct {
	DbName       string
	UserName     string
	UserPassword string
	Host         string
	Port         string
}

// Type for the db connection
type Pgconnector struct {
	Db *pg.DB
}

// Closes the underlying connection
func (p Pgconnector) Close() error {
	return p.Db.Close()
}

// Create a new Pgconnector for DbConfig
func CreateConnectClient(config DbConfig) Pgconnector {
	pgOptions := toPgOptions(config)

	// Connect to database
	db := pg.Connect(pgOptions)

	// Add logging longing for queries
	db.AddQueryHook(dbLogger{})
	return Pgconnector{
		Db: db,
	}
}

// Create DB tables if they dont exist
func (p Pgconnector) CreateSchemaIfNotExists() error {
	//Query pg_catalog.pg_database if a database with the name neo exists
	rows := []bool{}
	_, err := p.Db.Query(&rows, `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'neo');`)
	if err != nil {
		return err
	}

	//If the database is missing return an error
	if len(rows) == 0 || rows[0] == false {
		return errors.New(fmt.Sprintf("Database neo not found"))
	}

	//Create table if it isnt existing yet
	models := []interface{}{
		(*Neo)(nil),
	}
	for _, model := range models {
		fmt.Printf("Creating table for %s\n", reflect.TypeOf(model).Kind())
		p.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Convert DbConfig options to *pg.Options
func toPgOptions(config DbConfig) *pg.Options {
	return &pg.Options{
		Addr:            fmt.Sprintf("%s:%s", config.Host, config.Port),
		User:            config.UserName,
		Password:        config.UserPassword,
		Database:        config.DbName,
		ApplicationName: AppName,
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		PoolSize:        PoolSize,
		MinIdleConns:    MinIdleConns,
	}
}
