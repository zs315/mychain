package model

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)



/*
* 以下是http请求
*/

func MakeMuxRouter() http.Handler{
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", handleGetBlockchain).Methods("GET")
	r.HandleFunc("/", handleWriteBlock).Methods("POST")
	err:=http.ListenAndServe(":8888",r)
	if err!=nil {
		log.Fatalln("ListenAndServe err: ",err)
	}
	return r

}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	//log.Println("handleGetBlockchain begin.")
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))

	//log.Println("handleGetBlockchain end.")
}

type Message struct{
	BPM int
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	//log.Println("handleWriteBlock begin.")
	var m Message
	decoder := json.NewDecoder(r.Body)
	log.Println("decoder",decoder)
	if err := decoder.Decode(&m); err != nil {
		log.Println("handleWriteBlock err Decode: ",err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := GenerateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		log.Println("handleWriteBlock err generateBlock: ",err)
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		ReplaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

	//log.Println("handleWriteBlock begin.")
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
