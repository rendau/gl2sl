package http_api

import (
	"io/ioutil"
	"log"
	"net/http"
)

func hMsg(w http.ResponseWriter, r *http.Request) {
	cr := uReqGetCore(r)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	err = cr.HandleMessage(bodyBytes)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
}
