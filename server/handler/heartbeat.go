package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
)

func ReceivedHeartbeat(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	printBody(r)
	respondCode(w, http.StatusAccepted)
}

func LogHeartbeat(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	fmt.Println(r)
	respondCode(w, http.StatusAccepted)
}

func printBody(r *http.Request) {

	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		return
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Printf("BODY: %q", rdr1)
	r.Body = rdr2
}
