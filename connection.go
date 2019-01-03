/*
 *  Author MALDRU
 *  Email andres.latorre@ing-developers.com
 *  Copyright (c) 2019. All rights reserved.
 */
package go_mysql

import (
	"database/sql"
	// LIBRERIA PARA CONEXIONES MYSQL
	_ "github.com/go-sql-driver/mysql"
)

// ServerDB modelo para la conexion con la base de datos
type ServerDB struct {
	DSN      string `json:"dsn"` // si esta vacia se creara el dns por defecto con los datos siguientes
	Server   string `json:"server"`
	Port     string `json:"port"`
	DataBase string `json:"data_base"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Connect establece conexion con el motor de base de datos
func Connect(configServerDB ServerDB) (*MySql, error) {
	if configServerDB.DSN == "" {
		configServerDB.DSN = createDSN(configServerDB)
	}
	con, err := sql.Open("mysql", configServerDB.DSN)
	mySql := &MySql{
		db: con,
	}
	mySql.Connected = isConnected(mySql.db)
	return mySql, err
}

// createDSN Crea DSN por defecto
func createDSN(configServerDB ServerDB) string {
	return configServerDB.User + ":" + configServerDB.Password + "@tcp(" + configServerDB.Server + ":" + configServerDB.Port + ")/" + configServerDB.DataBase + "?parseTime=true&loc=America%2FBogota"
}

// isConnected verifica si la conexion se establecio correctamente
func isConnected(db *sql.DB) bool {
	return db.Ping() == nil
}
