package arangoproxy

import (
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

//DbConn is to establish connection with database.
func DbConn() driver.Database {

	dbUser := "user@localhost"
	dbPass := "password"
	dbName := "testdb"

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})

	if err != nil {
		panic(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(dbUser, dbPass),
	})
	if err != nil {
		panic(err)
	}

	db, err := client.Database(nil, dbName)
	if err != nil {
		panic(err)
	}

	log.Println("connecting to arangodb done")

	return db
}

//GetCollection is to get the collection from the databse.
func GetCollection(db driver.Database) driver.Collection {

	dbCollection := "tcpproxy"
	var col driver.Collection

	col, err := db.Collection(nil, dbCollection)
	if err != nil {
		panic(err)
	}

	return col
}

//CreateCollection is to create collection in database.
func CreateCollection(db driver.Database) driver.Collection {

	dbCollection := "tcpproxy"
	var col driver.Collection

	found, err := db.CollectionExists(nil, dbCollection)
	if err != nil {
		panic(err)
	}

	if found == true {
		col, err = db.Collection(nil, dbCollection)
		if err != nil {
			panic(err)
		}
	} else {
		options := &driver.CreateCollectionOptions{}
		col, err = db.CreateCollection(nil, dbCollection, options)
		if err != nil {
			panic(err)
		}
	}
	return col
}

// A PortIP is the struct for all the port, ip information.
type PortIP struct {
	Port int    `json:"port"`
	IP   string `json:"ip"`
}

// A ProxyServer is the struct for all the proxy information.
type ProxyServer struct {
	From     PortIP `json:"from"`
	To       PortIP `json:"to"`
	DataSize int64  `json:"datasize"`
	Active   bool   `json:"active"`
}

// A TotalServer structure is an interface with database
type TotalServer struct {
	ID string      `json:"id"`
	Ps ProxyServer `json:"ps"`
}
