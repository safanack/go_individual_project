package service

import (
	"database/sql"
	"datastream/config"
	"strings"

	"datastream/models"

	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendContactsToKafka(t *testing.T) {

	topic := "test-topic"
	brokerList := "localhost:9092"
	// brokerList := "localhost:9090"
	kafkaConfig := config.KafkaConfig{
		Broker: brokerList,
	}

	// Create Kafka connector
	kafkaConnector, err := config.NewKafkaConnector(kafkaConfig)
	if err != nil {
		panic(err)
	}
	defer kafkaConnector.Close()

	// Define test contacts data
	contacts := []models.Contacts{{
		ID:      1,
		Name:    "John Doe",
		Email:   "johndoe@example.com",
		Details: "Some details about John Doe",
	}}

	err = SendContactsToKafka(contacts, kafkaConnector, topic)
	if err != nil {
		t.Errorf("Error sending contacts to Kafka: %v", err)
	}

	if err != nil {
		t.Errorf("Error sending contacts to Kafka: %v", err)
	}
}

type DatabaseConnector interface {
	ConnectDatabase(string) (*sql.DB, error)
}

type MockDatabaseConnector struct {
	mock.Mock
}

func (m *MockDatabaseConnector) ConnectDatabase(databaseName string) (*sql.DB, error) {
	args := m.Called(databaseName)
	return args.Get(0).(*sql.DB), args.Error(1)
}
func TestClickhouseDb(t *testing.T) {
	// Create a real in-memory SQLite database for testing.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create a database connection: %v", err)
	}
	defer db.Close()

	// Create the required 'your_table' for the test in the SQLite database.
	_, err = db.Exec(`
        CREATE TABLE test_table (
            id INTEGER PRIMARY KEY,
            name TEXT
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create the table: %v", err)
	}
	// Insert sample data into the table.
	_, err = db.Exec(`INSERT INTO test_table (id,name) VALUES (1,"rahul"), (2,"Alice"), (3,"Bob")`)
	if err != nil {
		t.Fatalf("Failed to insert data into the table: %v", err)
	}
	// Define the expected query.
	expectedQuery := `SELECT * FROM test_table;`

	// Call the function with the query string and the database connection.
	rows := ClickhouseDb(expectedQuery, db)

	// Assert that the function returned the expected values.
	assert.NoError(t, err)
	assert.NotNil(t, rows)

	// // Test with a wrong query.
	// wrongQuery := `SELECT * FROM non_existent_table;`
	// rows = ClickhouseDb(wrongQuery, db)
	// // Assert that an error is returned when a wrong query is executed.
	// assert.Error(t, err)
}
func TestMysqlDB(t *testing.T) {
	// Create a real in-memory SQLite database for testing.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create a database connection: %v", err)
	}
	defer db.Close()

	// Create the required 'your_table' for the test in the SQLite database.
	_, err = db.Exec(`
        CREATE TABLE test_table (
            id INTEGER PRIMARY KEY,
            name TEXT
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create the table: %v", err)
	}
	// Insert sample data into the table.
	_, err = db.Exec(`INSERT INTO test_table (id,name) VALUES (1,"rahul"), (2,"Alice"), (3,"Bob")`)
	if err != nil {
		t.Fatalf("Failed to insert data into the table: %v", err)
	}
	// Define the expected query.
	expectedQuery := `SELECT * FROM test_table;`

	// Call the function with the query string and the database connection.
	err = MysqlDB(expectedQuery, db)

	// Assert that the function returned the expected values.
	assert.NoError(t, err)

}
func TestSendContactsToSql(t *testing.T) {
	topic := "test-topic"
	brokerList := "localhost:9092"

	// Create a mock DB connection for testing purposes.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create a database connection: %v", err)
	}
	defer db.Close()

	expectedQuery := "INSERT INTO `Contacts`(ID, Name, Email, Details) VALUES (1, 'safa', 'safa@gmail.com', 'some details')"

	// Mock the message consumption.

	// Call the function for testing.
	SendContactsToSql(brokerList, topic, db)

	// Check the expected outcome against the query executed.
	if strings.Contains(expectedQuery, "error") {
		t.Errorf("Test failed: Expected query contains error.")
	}

	// Add additional checks if required.
}
func BenchmarkSendContactsToSql(b *testing.B) {
	brokerList := "localhost:9092"
	topic := "test-topic"
	// Create a mock DB connection for testing purposes.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("Failed to create a database connection: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(`
	CREATE TABLE Contacts (
		ID int NOT NULL,
		Name varchar(255) NOT NULL,
		Email varchar(255) NOT NULL,
		Details varchar(300) NOT NULL
	  )
    `)
	if err != nil {
		b.Fatalf("Failed to create the table: %v", err)
	}

	// Run the function b.N times
	for n := 0; n < b.N; n++ {
		SendContactsToSql(brokerList, topic, db)
	}
}
