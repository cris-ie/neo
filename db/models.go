package db

import "time"

type Neo struct {
	Id                             string    `pg:",pk"`
	Name                           string    `pg:"name"`
	NasaJplUrl                     string    `pg:"jpl_url"`
	IsPotentiallyHazardousAsteroid bool      `pg:"hazardous,notnull,use_zero"`
	Date                           time.Time `pg:"date"`
}
