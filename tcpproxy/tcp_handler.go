package tcpproxy

import (
	"log"
	"time"

	ardb "github.com/TCP_Proxy/arangoproxy"
	"github.com/arangodb/go-driver"
)

// MapServer contains the mapping of proxy server id
// with the proxy server structure.
var MapServer map[string]TCPHandler

// MapData contains the mapping of proxy server id
// with data flowed through that proxy.
var MapData map[string]int64

// Startup function will start all the existing tcp proxies
// from the database.
func Startup() {
	MapData = make(map[string]int64)
	MapServer = make(map[string]TCPHandler)
	db := ardb.DbConn()
	col := ardb.CreateCollection(db)

	if col == nil {
		log.Println("collection not created")
	}

	query := "FOR d IN proxy RETURN d"
	cursor, err := db.Query(nil, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()

	ts := ardb.TotalServer{}

	for {
		meta, err := cursor.ReadDocument(nil, &ts.Ps)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		ts.ID = meta.Key

		if ts.Ps.Active == true {
			server := StartServer(ts)
			MapServer[ts.ID] = server
			//server.Tid <- ts.ID
		}
	}

	log.Println("startup completed")
}

// Calcdata is to caluclate data transfered in tcp proxy.
func Calcdata(seconds time.Duration) {
	db := ardb.DbConn()
	col := ardb.CreateCollection(db)
	query := "FOR d IN proxy RETURN d"

	for true {
		time.Sleep(seconds * time.Second)
		cursor, err := db.Query(nil, query, nil)
		if err != nil {
			panic(err)
		}

		ts := ardb.TotalServer{}

		for {
			meta, err := cursor.ReadDocument(nil, &ts.Ps)
			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				panic(err)
			}
			ts.ID = meta.Key
			total := MapData[ts.ID]
			patch := map[string]interface{}{
				"DataSize": total,
			}

			log.Printf("updating id %v data size %v", ts.ID, total)

			_, err = col.UpdateDocument(nil, ts.ID, patch)
			if err != nil {
				panic(err)
			}
		}

		cursor.Close()
	}
}

// Shutdown function is to close the tcp proxy
func (tcph *TCPHandler) Shutdown() {
	log.Printf("shutting down proxy port %s server %s \n", tcph.ts.Ps.From.IP, tcph.ts.Ps.To.IP)
	tcph.listen.Close()
	tcph.sender.Close()
}

// StartServer will start the tcp proxy server.
func StartServer(ts ardb.TotalServer) TCPHandler {
	var tcph TCPHandler

	tcph.ts = ts

	log.Printf("Lisetning on %s:%d Forwarding to %s:%d", tcph.ts.Ps.From.IP, tcph.ts.Ps.From.Port, tcph.ts.Ps.To.IP, tcph.ts.Ps.To.Port)
	tcplns(&tcph)

	return tcph
}
