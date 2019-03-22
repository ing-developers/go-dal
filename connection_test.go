/*
 * Author MALDRU
 * Email andres.latorre@ing-developers.com
 * Copyright (c) 2019. All rights reserved.
 */

package go_mysql

import (
	"github.com/ing-developers/go-tools"
	"testing"
	// LIBRERIA PARA CONEXIONES MYSQL
	_ "github.com/go-sql-driver/mysql"
)

func TestIsConnected(t *testing.T) {
	var configServerDB ServerDB
	err := tools.Decode("./serverDB.json", &configServerDB)
	if err != nil {
		t.Error(err)
		return
	}
	mysql, err := Connect(configServerDB)
	if err != nil {
		t.Error(err)
		return
	}
	if mysql.Connected {
		t.Log("Se establecio coneccion")
	} else {
		t.Log("No se establecio coneccion")
	}
}
