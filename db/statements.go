package db

import "fmt"

func QueryDoesDbExist(dbName string) string {
	return fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');`, dbName)
}
