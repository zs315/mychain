package main

import (
	"github.com/mychain/model"
	"time"
	"github.com/davecgh/go-spew/spew"
)


func main(){
	t := time.Now()
	genesisBlock := model.Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	model.Blockchain = append(model.Blockchain, genesisBlock)
	model.MakeMuxRouter()
}