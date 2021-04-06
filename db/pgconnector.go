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

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	fmt.Println(string(fq))
	return nil
}

const (
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
	PoolSize     = 10
	MinIdleConns = 10
	AppName      = "neo-app"
)

type DbConfig struct {
	DbName       string
	UserName     string
	UserPassword string
	Host         string
	Port         string
}

type Pgconnector struct {
	Db *pg.DB
}

func (p Pgconnector) GetConnection() *pg.DB {
	return p.Db
}

func (p Pgconnector) Close() error {
	return p.Db.Close()
}

func toPgOptions(config DbConfig) *pg.Options {
	return &pg.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),

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
func CreateConnectClient(config DbConfig) Pgconnector {
	pgOptions := toPgOptions(config)
	db := pg.Connect(pgOptions)
	db.AddQueryHook(dbLogger{})
	return Pgconnector{
		Db: db,
	}
}
func (p Pgconnector) CreateSchemaIfNotExists() error {
	rows := []bool{}
	_, err := p.Db.Query(&rows, QueryDoesDbExist("neo"))
	if err != nil {
		return err
	}

	if len(rows) == 0 || !rows[0] {
		return errors.New(fmt.Sprintf("Database neo not found"))
	}
	//No schema -> create
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
