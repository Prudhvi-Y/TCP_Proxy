package interfacedb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	ardb "github.com/TCP_Proxy/arangoproxy"
	tcp "github.com/TCP_Proxy/tcpproxy"
	driver "github.com/arangodb/go-driver"
	"github.com/gorilla/mux"
)

//ShowAll functions shows all the proxies in database
func ShowAll(w http.ResponseWriter, r *http.Request) {
	log.Println("executing ShowAll")

	if !authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	db := ardb.DbConn()

	query := "FOR d IN tcpproxy RETURN d"
	cursor, err := db.Query(nil, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()

	ts := ardb.TotalServer{}
	ss := []ardb.TotalServer{}

	for {
		meta, err := cursor.ReadDocument(nil, &ts.Ps)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		ts.ID = meta.Key
		ts.Ps.DataSize = tcp.MapData[ts.ID]
		ss = append(ss, ts)
	}

	if err := json.NewEncoder(w).Encode(ss); err != nil {
		panic(err)
	}
}

//Show function shows the proxy details of given id
func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["showId"]

	log.Println("executing Show", key)
	if !authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	db := ardb.DbConn()
	col := ardb.GetCollection(db)

	ps, err := utilreaddoc(col, key)
	if err != nil {
		panic(err)
	}

	ps.Ps.DataSize = tcp.MapData[ps.ID]

	log.Println(ps)

	if err := json.NewEncoder(w).Encode(ps); err != nil {
		panic(err)
	}
}

//Insert function will insert a new proxy to database
func Insert(w http.ResponseWriter, r *http.Request) {

	log.Println("executing Insert")
	if !authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		db := ardb.DbConn()
		col := ardb.GetCollection(db)

		var ps ardb.ProxyServer

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &ps); err != nil {
			w.Header().Set("Content-Type", "application/json; chatset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		if err := validatedata(ps); err != nil {
			panic(err)
		}

		meta, err := col.CreateDocument(nil, ps)
		if err != nil {
			panic(err)
		}

		rs, err := utilreaddoc(col, meta.Key)
		if err != nil {
			panic(err)
		}
		rs.ID = meta.Key

		if rs.Ps.Active == true {
			server := tcp.StartServer(rs)
			tcp.MapServer[rs.ID] = server
			//server.Tid <- rs.ID
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(rs); err != nil {
			panic(err)
		}

	}
}

//Delete will shutdown and remove the proxy of given id from the database.
func Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["deleteId"]

	log.Println("executing Delere", key)
	if !authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method == "POST" {
		db := ardb.DbConn()
		col := ardb.GetCollection(db)

		rs, err := utilreaddoc(col, key)
		if err != nil {
			panic(err)
		}

		if _, ok := tcp.MapServer[rs.ID]; ok {
			server := tcp.MapServer[rs.ID]
			//log.Println(server.Server.Addr)
			server.Shutdown()
			delete(tcp.MapServer, rs.ID)
			delete(tcp.MapData, rs.ID)
		}

		_, err = col.RemoveDocument(nil, key)
		if err != nil {
			panic(err.Error())
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusFound)
		if err := json.NewEncoder(w).Encode(rs); err != nil {
			panic(err)
		}
	}
}

//Update the details of given proxy id
func Update(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["updateId"]

	log.Println("executing Update", key)
	if !authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method == "POST" {
		db := ardb.DbConn()
		col := ardb.GetCollection(db)

		var ps ardb.ProxyServer

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &ps); err != nil {
			w.Header().Set("Content-Type", "application/json; chatset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		log.Println(ps)

		if err := validateupdatedata(ps); err != nil {
			panic(err)
		}

		patch := createpatch(ps)

		_, err = col.UpdateDocument(nil, key, patch)
		if err != nil {
			panic(err)
		}

		log.Println("updated document")

		rs, err := utilreaddoc(col, key)
		if err != nil {
			panic(err)
		}

		if _, ok := tcp.MapServer[rs.ID]; ok {
			server := tcp.MapServer[rs.ID]
			//log.Println(server.Server.Addr)
			server.Shutdown()
			delete(tcp.MapServer, rs.ID)
			delete(tcp.MapData, rs.ID)
		}

		if rs.Ps.Active == true {
			server := tcp.StartServer(rs)
			tcp.MapServer[rs.ID] = server
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusFound)
		if err := json.NewEncoder(w).Encode(rs); err != nil {
			panic(err)
		}
	}
}
