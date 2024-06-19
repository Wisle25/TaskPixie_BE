package services

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/lib/pq"
	"github.com/wisle25/be-template/commons"
)

// ConnectDB initializes a connection to the PostgreSQL cache using the provided configuration.
func ConnectDB(config *commons.Config) *sql.DB {
	var err error
	dbName := config.DBName

	// Use the test cache if the application environment is set to 'dev'.
	if config.AppEnv == "test" {
		dbName = config.DBNameTest
	}

	// Format the data source name (DSN) string for connecting to the cache.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		dbName,
		config.DBPort,
	)

	// Open a connection to the cache.
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("connect_DB_err: Opening: %v", err))
	}

	// Ping the cache to ensure the connection is established.
	if err = DB.Ping(); err != nil {
		panic(fmt.Errorf("connect_DB_err: Pinging: %v", err))
	}

	log.Println("Successfully connected to Postgres!")

	return DB
}

// GetTableDB parses SQL query results into a slice of the specified type T.
func GetTableDB[T any](rows *sql.Rows) []T {
	var table []T

	// Iterate through the rows returned by the query.
	for rows.Next() {
		var data T

		s := reflect.ValueOf(&data).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)

		// Prepare the columns slice with pointers to the fields of the struct.
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		// Scan the row into the struct fields.
		if err := rows.Scan(columns...); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil // Return nil if there is an error
		}

		// Append the struct to the table slice.
		table = append(table, data)
		log.Printf("Scanned row: %+v", data)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows: %v", err)
		return nil // Return nil if there is an error
	}

	return table
}
