/*
 * Author MALDRU
 * Email andres.latorre@ing-developers.com
 * Copyright (c) 2019. All rights reserved.
 */

package go_dal

import (
	"database/sql"
	"encoding/json"
)

// DAL modelo para compartir la conexion con el servidor de base de datos entre las funciones existentes
type DAL struct {
	db           *sql.DB
	tx           *sql.Tx
	stmt         *sql.Stmt
	Connected    bool
	LastID       int64
	AffectedRows int64
	Errors       errors
}

//Rows estructura para mapeo de datos obtenidos
type Rows map[string]string

//errors Lista de errores generados
type errors []error

//CloseConnection cierra la conexion actual
func (m *DAL) CloseConnection() error {
	return m.db.Close()
}

// isTransaction verifica si es una transaccion
func (m *DAL) isTransaction() bool {
	return m.tx != nil
}

// isSTMT verifica si es una sentencia preparada
func (m *DAL) isSTMT() bool {
	return m.stmt != nil
}

// BeginTransaction inicia transaccion
func (m *DAL) BeginTransaction() (err error) {
	m.tx, err = m.db.Begin()
	return err
}

// GetRowsSTMT obtiene las filas de una sentencia preparada, la primera vez a ejecutarse se necesitara el query
func (m *DAL) GetRowsSTMT(query string, values ...interface{}) (table []Rows, err error) {
	var rows *sql.Rows

	if m.stmt == nil || query != "" {
		if m.isTransaction() {
			m.stmt, err = m.tx.Prepare(query)
		} else {
			m.stmt, err = m.db.Prepare(query)
		}
	}

	rows, err = m.stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}

// GetRowsQuery obtiene las filas de una consulta SQL
func (m *DAL) GetRowsQuery(query string, values ...interface{}) (table []Rows, err error) {
	var rows *sql.Rows

	if m.isTransaction() {
		rows, err = m.tx.Query(query, values...)
	} else {
		rows, err = m.db.Query(query, values...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}

// ExecuteSTMT ejecuta una sentencia preparada, la primera vez a ejecutarse se necesitara el query
func (m *DAL) ExecuteSTMT(query string, values ...interface{}) (err error) {
	var result sql.Result

	if m.stmt == nil || query != "" {
		if m.isTransaction() {
			m.stmt, err = m.tx.Prepare(query)
		} else {
			m.stmt, err = m.db.Prepare(query)
		}
	}

	result, err = m.stmt.Exec(values...)
	if err == nil {
		m.LastID, err = result.LastInsertId()
		m.AffectedRows, err = result.RowsAffected()
	} else {
		m.Errors = append(m.Errors, err)
	}
	return err
}

// ExecuteQuery ejecuta una consulta sql
func (m *DAL) ExecuteQuery(query string, values ...interface{}) (err error) {
	var result sql.Result

	if m.isTransaction() {
		result, err = m.tx.Exec(query, values...)
	} else {
		result, err = m.db.Exec(query, values...)
	}

	if err == nil {
		m.LastID, err = result.LastInsertId()
		m.AffectedRows, err = result.RowsAffected()
	} else {
		m.Errors = append(m.Errors, err)
	}
	return err
}

// ToStruct convierte un maps a un struct
func ToStruct(row Rows, model interface{}) error {
	js, err := json.Marshal(row)
	if err == nil {
		err = json.Unmarshal(js, &model)
	}
	return err
}

// ToSliceOfStructs convierte un maps a un slice de structs
func ToSliceOfStructs(table []Rows, model interface{}) error {
	js, err := json.Marshal(table)
	if err == nil {
		err = json.Unmarshal(js, &model)
	}
	return err
}

// FinalizeTransaction finaliza la transaccion activa verificando errores y realiza rollback si es el caso
func (m *DAL) FinalizeTransaction() (err error) {
	if len(m.Errors) > 0 {
		err = m.tx.Rollback()
	} else {
		err = m.tx.Commit()
	}
	return err
}

// rowsToMap genera maps de resultado de consulta sql
func rowsToMap(rows *sql.Rows) (table []Rows, err error) {
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for c := range values {
		scanArgs[c] = &values[c]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		row := make(Rows, len(columns))
		for i, val := range values {
			if val == nil {
				row[columns[i]] = "--"
			} else {
				row[columns[i]] = string(val)
			}
		}
		table = append(table, row)
	}
	return table, err
}
