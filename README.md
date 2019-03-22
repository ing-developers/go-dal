# go_mysql
Conexión y operaciones con el motor de base de datos MySql o MariaDB con **Go**
## 1. Instalacion
Para la instalacion de **go_mysql** ejecutar el siguiente comando
```bash
go get -u github.com/ing-developers/go-dal
```
## Instalar Dependencias
Las dependencias requeridas para utilizar  **go_mysql** son las siguientes:
### Utilizar el bat ubicado en el directorio go_mysql
1. Dirigirse al directorio de **go_mysql**
```bash
${GOPATH}\src\github.com\ing-developers\go-mysql`
```
2. Dentro del directorio debe ejecutar el archivo **package.bat**

### Instalar manualmente las dependencias
1. https://github.com/go-sql-driver/mysql
```bash
go get -u github.com/go-sql-driver/mysql
```
2. Opcional para establecer conexion con archivo JSON https://github.com/ing-developers/go-tools
```bash
go get -u github.com/ing-developers/go-tools
```

## 2. Uso
## Coneccion con el servidor mysql o mariaDB
En donde se deseé establecer la conexion con el motor de base de datos, importar la libreria que manejara dicha conexion.

Por ejemplo, con mysql o mariaDB

```go
import (	
	// LIBRERIA PARA CONEXIONES MYSQL
	_ "github.com/go-sql-driver/mysql"
)
```

- **Mediante la struct definida**
```go
configServerDB := go_mysql.ServerDB{	
	Server:   "localhost",
	Port:     "3306",
	DataBase: "test",
	User:     "root",
	Password: "123321",
}
mysql, err := Connect(configServerDB)
if err != nil {
	log.Fatal(err)
	return
}
if mysql.Connected {
	log.Println("Se establecio coneccion")
} else {
	log.Println("No se establecio coneccion")
}
```
**o especificando el DSN (Data Source Name)**
```go
configServerDB := go_mysql.ServerDB{	
	DSN:   "usuario:contraseña@tcp(host:puerto)/nombreBD",
}
...
```
- **Mediante archivo JSON**
**Archivo JSON**
```json
{
	"server": "localhost",
	"port": "3306",
	"data_base": "test",
	"user": "root",
	"password": "123456"
}
```
**o especificando el DSN (Data Source Name)**
```json
{
	"dsn": "usuario:contraseña@tcp(host:puerto)/nombreBD"
}
```
**Archivo de conexion**
```go
var configServerDB go_mysql.ServerDB
err := tools.Decode("./serverDB.json", &configServerDB)
if err != nil {
	log.Fatal(err)
	return
}
mysql, err := Connect(configServerDB)
if err != nil {
	log.Fatal(err)
	return
}
if mysql.Connected {
	log.Println("Se establecio coneccion")
} else {
	log.Println("No se establecio coneccion")
}
```

## Obtener filas
```go
...
if mysql.Connected {
	table, err := mysql.GetRowsQuery("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}
	table, err = mysql.GetRowsQuery("SELECT * FROM products WHERE id = ?", 3)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}
}
...
```
**Mediante sentencia preparada**
```go
...
if mysql.Connected {
	table, err := mysql.GetRowsSTMT("SELECT * FROM products WHERE id = ?", 1)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}

	table, err = mysql.GetRowsSTMT("", 2)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}

	table, err = mysql.GetRowsSTMT("SELECT * FROM products WHERE name = ?", "product")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}

	table, err = mysql.GetRowsSTMT("", "product 2")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(table)
	}
}
...
```
## Conversion a struct o slice de struct
**Slice de struct**
```go
table, err := mysql.GetRowsQuery("SELECT * FROM products")
if err != nil {
	log.Fatal(err)
} else {
	type product struct {
		ID   int64  `json:"id,string"`
		Name string `json:"name"`
	}
	var products []product
	err = go_mysql.ToSliceOfStructs(table, &products)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(products)
}
```
**Struct**
```go
table, err := mysql.GetRowsQuery("SELECT * FROM products WHERE id = ?", 1)
if err != nil {
	log.Fatal(err)
} else {
	type product struct {
		ID   int64  `json:"id,string"`
		Name string `json:"name"`
	}
	var productOne product
	err = go_mysql.ToStruct(table[0], &productOne)
	if err != nil {
	log.Fatal(err)
	return
	}
	log.Println("Id: ", productOne.ID, " Name: ", productOne.Name)
}
```
## Ejecutar consulta SQL
```go
...
if mysql.Connected {
	err := mysql.ExecuteQuery("INSERT INTO products VALUES (null, ?)", "product 9")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Insercion correcta, id = ", mysql.LastID)
	}
}
```
**Mediante sentencia preparada**
```go
...
if mysql.Connected {
	err := mysql.ExecuteSTMT("INSERT INTO products VALUES (null, ?)", "product 4")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Insercion correcta, id = ", mysql.LastID)
	}
	err = mysql.ExecuteSTMT("", "product 5")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Insercion correcta, id = ", mysql.LastID)
	}
	err = mysql.ExecuteSTMT("UPDATE products SET name = ? WHERE id = ?", "product five", 5)
	if err != nil {
		log.Fatal(err)
	} else {
	log.Println("Actualizacion correcta, filas afectadas: ", mysql.AffectedRows)
	}
	err = mysql.ExecuteSTMT("", "product four", 4)
	if err != nil {
		log.Fatal(err)
	} else {
	log.Println("Actualizacion correcta, filas afectadas: ", mysql.AffectedRows)
	}
}
...
```
## Transacciones
```go
...
if mysql.Connected {
	err := mysql.BeginTransaction()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = mysql.ExecuteSTMT("UPDATE products SET name = ? WHERE id = ?", "product One", 1)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Actualizacion correcta, filas afectadas: ", mysql.AffectedRows)
	}
	err = mysql.ExecuteSTMT("", "product Two", 2)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Actualizacion correcta, filas afectadas: ", mysql.AffectedRows)
	}
	err = mysql.ExecuteQuery("INSERT INTO products VALUES (null, ?)", "product 6")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Insercion correcta, id = ", mysql.LastID)
	}

	wasCommit, err := mysql.FinalizeTransaction()

	if err != nil {
		if wasCommit {
			log.Println("Error al ejecutar el commit")
		} else {
			log.Println("Error al ejecutar el rollback")
		}
	} else {
		if wasCommit {
			log.Println("Commit correcto")
		} else {
			log.Println("rollback correcto")
		}
	}
}
```
