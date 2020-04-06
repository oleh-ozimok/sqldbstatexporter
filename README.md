# Golang SQL DB Stat Exporter
Prometheus exporter for Golang sql DB stats

## Usage

```go

import (
    "database/sql"

    "github.com/oleh-ozimok/sqldbstatexporter"
    "github.com/prometheus/client_golang/prometheus"
    
    _ "github.com/go-sql-driver/mysql"
)

db, _ := sql.Open("mysql", "user:password@/dbname")

prometheus.Register(
    sqldbstatxporter.New(db, "mysql", prometheus.Labels{
        "pool_id": "main",
    }), 	
)




    

```