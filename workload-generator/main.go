package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	_ "github.com/go-sql-driver/mysql"   // MySQL driver
	_ "github.com/lib/pq"                // PostgreSQL driver
)

// Config struct to hold the configuration
type Config struct {
	DB_platform            string `json:"db_platform"`
	Server                 string `json:"server"`
	Database               string `json:"database"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	Port                   int    `json:"port"`
	MAX_PARALLEL_WORKLOADS int32  `json:"MAX_PARALLEL_WORKLOADS"`
}

const (
	configFileName                = "config.json"
	createDatabaseScriptSQLServer = "CREATE_DATABASE.sql"
	createTablesScriptSQLServer   = "CREATE_TABLES.sql"
	createWorkloadSprocSQLServer  = "CREATE_WORKLOAD_SPROC.sql"
	createDatabaseScriptPostgres  = "CREATE_DATABASE_PG.sql"
	createTablesScriptPostgres    = "CREATE_TABLES_PG.sql"
	createWorkloadSprocPostgres   = "CREATE_WORKLOAD_SPROC_PG.sql"
	createDatabaseScriptMySQL     = "CREATE_DATABASE_MYSQL.sql"
	createTablesScriptMySQL       = "CREATE_TABLES_MYSQL.sql"
	createWorkloadSprocMySQL      = "CREATE_WORKLOAD_SPROC_MYSQL.sql"
)

var (
	createDatabaseSQL string
	createTablesSQL   string
	createWorkloadSQL string
	workloadProcCall  string
	driverName        string
)

func main() {
	// Read the configuration file
	configFile, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Parse the JSON configuration
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Determine the database platform
	switch config.DB_platform {
	case "sqlserver":
		createDatabaseSQL = createDatabaseScriptSQLServer
		createTablesSQL = createTablesScriptSQLServer
		createWorkloadSQL = createWorkloadSprocSQLServer
		workloadProcCall = "exec dbo.demo_workload"
		driverName = "sqlserver"
		// Create the connection string for master database
		connString_master := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master", config.Username, config.Password, config.Server, config.Port)
		// Create the connection string for the target database
		connString_db := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", config.Username, config.Password, config.Server, config.Port, config.Database)
		createDatabase(connString_master)
		createTables(connString_db)
		createWorkloadSproc(connString_db)
		startWorkload(connString_db, config.MAX_PARALLEL_WORKLOADS)
	case "postgres":
		createDatabaseSQL = createDatabaseScriptPostgres
		createTablesSQL = createTablesScriptPostgres
		createWorkloadSQL = createWorkloadSprocPostgres
		workloadProcCall = "SELECT demo_workload.demo_workload()"
		driverName = "postgres"
		// Create the connection string for the target database
		connString_db_pg := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.Server, config.Port, config.Database)
		createDatabase(connString_db_pg)
		createTables(connString_db_pg)
		createWorkloadSproc(connString_db_pg)
		startWorkload(connString_db_pg, config.MAX_PARALLEL_WORKLOADS)
	case "mysql":
		createDatabaseSQL = createDatabaseScriptMySQL
		createTablesSQL = createTablesScriptMySQL
		createWorkloadSQL = createWorkloadSprocMySQL
		driverName = "mysql"
		// Create the connection string for the target database
		connString_db_mysql := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Server, config.Port, config.Database)
		createDatabase(connString_db_mysql)
		createTables(connString_db_mysql)
		createWorkloadSproc(connString_db_mysql)
		startWorkload(connString_db_mysql, config.MAX_PARALLEL_WORKLOADS)
	default:
		log.Fatalf("NO CONFIG SPECIFIED. EXITING...")
	}
}

func startWorkload(connstring string, MAX_PARALLEL_WORKLOADS int32) {
	// Open a connection to the database
	db := openDatabaseConnection(connstring)

	// Run the workload
	var counter int64 = 0
	var current_workloads int32 = 0
	for {

		//fmt.Println("current_workloads: ", atomic.LoadInt32(&current_workloads))
		if atomic.LoadInt32(&current_workloads) < MAX_PARALLEL_WORKLOADS {
			counter++
			atomic.AddInt32(&current_workloads, 1)
			go runWorkload(db, &current_workloads, counter)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func createWorkloadSproc(connstring string) {
	// Open a connection to the database
	db := openDatabaseConnection(connstring)

	fmt.Println("Creating the workload stored proc...")
	// Read the SQL query from an external file
	query, err := os.ReadFile(createWorkloadSQL)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	// Execute the query
	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatalf("Error running query to create stored procs: %v", err)
	}

	fmt.Println("setup complete...")
}

func createDatabase(connstring string) {
	// Open a connection to the database
	db := openDatabaseConnection(connstring)

	// Define the SQL query
	fmt.Println("Setting up the workload database...")
	// Read the SQL query from an external file
	query, err := os.ReadFile(createDatabaseSQL)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	// Execute the query
	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatalf("Error running query to create database: %v", err)
	}
}

func createTables(connstring string) {
	// Open a connection to the database
	db := openDatabaseConnection(connstring)

	// Define the SQL query
	fmt.Println("Setting up the workload tables...")
	// Read the SQL query from an external file
	query, err := os.ReadFile(createTablesSQL)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	// Execute the query
	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatalf("Error running query to create tables: %v", err)
	}
}

func runWorkload(db *sql.DB, current_workloads *int32, workloadcounter int64) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	//create a random integer between 1 and 1000
	sleepTime := rand.Intn(1000) + 1
	fmt.Println("Sleeping for ", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	// Call the stored procedure
	fmt.Printf("Running workload iteration %d...\n", workloadcounter)
	fmt.Println("current_workloads: ", atomic.LoadInt32(current_workloads))

	_, err := db.Exec(fmt.Sprintf(workloadProcCall))
	if err != nil {
		fmt.Println("Error calling stored procedure: %v . Waiting 5 seconds before calling next proc.. ", err)
		time.Sleep(5 * time.Second)
		return
	}

	fmt.Println("Stored procedure executed successfully! for ", workloadcounter)
	atomic.AddInt32(current_workloads, -1)
}

func openDatabaseConnection(connstring string) *sql.DB {
	// Open a connection to the database
	db, err := sql.Open(driverName, connstring)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	//verify connection
	// Retry pinging the database until successful
	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fmt.Println("Unable to connect to the database. Retrying in 5 seconds... %v", err)
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Connected to the database successfully!")

	return db
}
