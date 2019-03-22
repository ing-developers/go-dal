/*
 *  Author MALDRU
 *  Email andres.latorre@ing-developers.com
 *  Copyright (c) 2019. All rights reserved.
 */
package go_dal

import (
	"database/sql"
	"fmt"
	"log"
)

// ServerDB modelo para la conexion con la base de datos
type ServerDB struct {
	Engine   string `json:"engine"`
	DSN      string `json:"dsn"` // si esta vacia se creara el dns por defecto con los datos siguientes
	Server   string `json:"server"`
	Port     string `json:"port"`
	DataBase string `json:"data_base"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Connect establece conexion con el motor de base de datos
func Connect(configServerDB ServerDB) (*DAL, error) {
	if configServerDB.DSN == "" {
		configServerDB.DSN = createDSN(configServerDB)
	}
	con, err := sql.Open(configServerDB.Engine, configServerDB.DSN)
	engine := &DAL{
		db: con,
	}
	engine.Connected = isConnected(engine.db)
	return engine, err
}

// createDSN Crea DSN por defecto
func createDSN(configServerDB ServerDB) string {
	switch configServerDB.Engine {
	case "mysql":
		return configServerDB.User + ":" + configServerDB.Password + "@tcp(" + configServerDB.Server + ":" + configServerDB.Port + ")/" + configServerDB.DataBase + "?parseTime=true&loc=America%2FBogota"
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", configServerDB.User, configServerDB.Password, configServerDB.Server, configServerDB.Port, configServerDB.DataBase)
	default:
		log.Fatal("Engine no supported")
		return ""
	}
}

// isConnected verifica si la conexion se establecio correctamente
func isConnected(db *sql.DB) bool {
	return db.Ping() == nil
}
