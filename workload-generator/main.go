package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"time"
	"io/ioutil"
	"math/rand"
	"encoding/json"

	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

// Config struct to hold the configuration
type Config struct {
    Server   string `json:"server"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
    Port     int    `json:"port"`
	MAX_PARALLEL_WORKLOADS int32 `json:"MAX_PARALLEL_WORKLOADS"`
}

func main() {
	// Read the configuration file
    configFile, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    // Parse the JSON configuration
    var config Config
    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatalf("Error parsing config file: %v", err)
    }

    // Create the connection string for master database
    connString_master := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master",
        config.Username, config.Password, config.Server, config.Port)

		// Create the connection string for the target database
	connString_db := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
	config.Username, config.Password, config.Server, config.Port, config.Database)
 
	createDatabase(connString_master)

	createTables(connString_db)
	
	createWorkloadSproc(connString_db)
	
	startWorkload(connString_db, config.MAX_PARALLEL_WORKLOADS)
}

func startWorkload(connString string, MAX_PARALLEL_WORKLOADS int32) {
	// Open a connection to the target database
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer db.Close()

	// Retry pinging the target database until successful
	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fmt.Printf("Unable to connect to the target database. Retrying in 5 seconds... %v\n", err)
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Connected to the target database successfully!")

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

func createWorkloadSproc(connString string) {
	// Open a connection to the target database
    db, err := sql.Open("sqlserver", connString)
    if err != nil {
        log.Fatalf("Error creating connection pool: %v", err)
    }
    defer db.Close()

    // Retry pinging the target database until successful
    for {
        err = db.Ping()
        if err == nil {
            break
        }
        fmt.Printf("Unable to connect to the target database. Retrying in 5 seconds... %v\n", err)
        time.Sleep(5 * time.Second)
    }
    fmt.Println("Connected to the target database successfully!")

	fmt.Println("Creating the workload stored proc...")	
	// Read the SQL query from an external file
    query, err := ioutil.ReadFile("CREATE_WORKLOAD_SPROC.sql")
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

func createDatabase(connString string) {
	// Open a connection to the database
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer db.Close()

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

	// Define the SQL query
	fmt.Println("Setting up the workload database...")	
	// Read the SQL query from an external file
    query, err := ioutil.ReadFile("CREATE_DATABASE.sql")
    if err != nil {
        log.Fatalf("Error reading SQL file: %v", err)
    }

    // Execute the query
    _, err = db.Exec(string(query))
    if err != nil {
        log.Fatalf("Error running query to create database: %v", err)
    }
}

func createTables(connString string) {
	// Open a connection to the database
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer db.Close()

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

	// Define the SQL query
	fmt.Println("Setting up the workload tables...")	
	// Read the SQL query from an external file
    query, err := ioutil.ReadFile("CREATE_TABLES.sql")
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
	procName := "dbo.demo_workload"
	_, err := db.Exec(fmt.Sprintf("EXEC %s", procName))
	if err != nil {
		fmt.Println("Error calling stored procedure: %v . Waiting 5 seconds before calling next proc.. ", err)
		time.Sleep(5 * time.Second)
		return
	}

	fmt.Println("Stored procedure executed successfully! for ", workloadcounter)
	atomic.AddInt32(current_workloads, -1)	
}
