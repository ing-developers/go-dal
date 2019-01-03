/*
 * Author MALDRU
 * Email andres.latorre@ing-developers.com
 * Copyright (c) 2019. All rights reserved.
 */

package go_mysql

import (
	"github.com/ing-developers/go-tools"
	"testing"
)

func TestMySql_GetRowsQuery(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		table, err := mysql.GetRowsQuery("SELECT * FROM products")
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}
		table, err = mysql.GetRowsQuery("SELECT * FROM products WHERE id = ?", 3)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestMySql_ExecuteQuery(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		err := mysql.ExecuteQuery("INSERT INTO products VALUES (null, ?)", "product 2")
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestMySql_GetRowsSTMT(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		table, err := mysql.GetRowsSTMT("SELECT * FROM products WHERE id = ?", 1)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}

		table, err = mysql.GetRowsSTMT("", 2)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}

		table, err = mysql.GetRowsSTMT("SELECT * FROM products WHERE name = ?", "product")
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}

		table, err = mysql.GetRowsSTMT("", "product 2")
		if err != nil {
			t.Error(err)
		} else {
			t.Log(table)
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestMySql_ExecuteSTMT(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		err := mysql.ExecuteSTMT("INSERT INTO products VALUES (null, ?)", "product 4")
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}
		err = mysql.ExecuteSTMT("", "product 5")
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}
		err = mysql.ExecuteSTMT("UPDATE products SET name = ? WHERE id = ?", "product five", 5)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}
		err = mysql.ExecuteSTMT("", "product four", 4)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}

	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestMySql_BeginTransaction(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		err := mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
			return
		}
		err = mysql.ExecuteSTMT("UPDATE products SET name = ? WHERE id = ?", "product One", 1)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Actualizacion correcta, id = ", mysql.LastID)
		}
		err = mysql.ExecuteSTMT("", "product Two", 2)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}
		err = mysql.ExecuteQuery("INSERT INTO products VALUES (null, ?)", "product 6")
		if err != nil {
			t.Error(err)
		} else {
			t.Log("Insercion correcta, id = ", mysql.LastID)
		}

		wasCommit, err := mysql.FinalizeTransaction()

		if err != nil {
			if wasCommit {
				t.Log("Error al ejecutar el commit")
			} else {
				t.Log("Error al ejecutar el rollback")
			}
		} else {
			if wasCommit {
				t.Log("Commit correcto")
			} else {
				t.Log("rollback correcto")
			}
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestToSliceOfStructs(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		table, err := mysql.GetRowsQuery("SELECT * FROM products")
		if err != nil {
			t.Error(err)
		} else {
			type product struct {
				ID   int64  `json:"id,string"`
				Name string `json:"name"`
			}
			var products []product
			err = ToSliceOfStructs(table, &products)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(products)
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}

func TestToStruct(t *testing.T) {
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
	defer mysql.CloseConnection()
	if mysql.Connected {
		t.Log("Se establecio coneccion")
		table, err := mysql.GetRowsQuery("CALL getByID(?)", 1)
		if err != nil {
			t.Error(err)
		} else {
			type product struct {
				ID   int64  `json:"id,string"`
				Name string `json:"name"`
			}
			var productOne product
			err = ToStruct(table[0], &productOne)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log("Id: ", productOne.ID, " Name: ", productOne.Name)
		}
	} else {
		t.Log("No se establecio coneccion")
	}
}
