package interfacedb

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	ardb "github.com/TCP_Proxy/arangoproxy"
	driver "github.com/arangodb/go-driver"
)

func validatedata(ps ardb.ProxyServer) error {
	if net.ParseIP(ps.From.IP) != nil {
		//if net.ParseIP(ps.ToIP) != nil {
		if ps.From.Port != 0 { //&& doc.Ps.ToPORT != "" && doc.Ps.DataSize > 0 {
			return nil
		}
		//}
	}
	log.Println(ps)

	return errors.New("incorrect data format")
}

func validateupdatedata(ps ardb.ProxyServer) error {

	if ps.From.IP != "" {
		if net.ParseIP(ps.From.IP) == nil {
			return errors.New("incorrect data format")
		}
	}

	log.Println(ps)

	return nil
}

func createpatch(ps ardb.ProxyServer) interface{} {
	patch := make(map[string]interface{})

	if ps.From.IP != "" {
		patch["fromip"] = ps.From.IP
	}

	if ps.From.Port != 0 {
		patch["fromport"] = ps.From.Port
	}

	patch["toport"] = ps.To.Port

	if ps.To.IP != "" {
		patch["toip"] = ps.To.IP
	}

	patch["active"] = ps.Active

	if ps.DataSize == 0 {
		patch["datasize"] = ps.DataSize
	}

	return patch
}

func utilreaddoc(col driver.Collection, key string) (ardb.TotalServer, error) {
	var ts ardb.TotalServer

	_, err := col.ReadDocument(nil, key, &ts.Ps)
	if err != nil {
		panic(err)
	}
	ts.ID = key

	return ts, err
}

func authenticate(r *http.Request) bool {
	s := r.Header.Get("Authorization")
	pwd, _ := os.Getwd()
	authn, err := ioutil.ReadFile(pwd + "/interfacedb/file.txt")
	if err != nil {
		log.Println("Error loading file")
	}
	auth := string(authn)
	if s[7:] == auth {
		return true
	}
	return false
}
