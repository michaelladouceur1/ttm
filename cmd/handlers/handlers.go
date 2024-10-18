package handlers

import (
	"fmt"
	"os"
	ttmConfig "ttm/pkg/config"
	ttmStore "ttm/pkg/store"
	ttmDB "ttm/pkg/store/db"
)

var config *ttmConfig.Config
var store *ttmStore.Store

func init() {
	store = ttmStore.NewStore(ttmDB.NewDBStore())
	store.Init()

	ttmConfig.Init()
	conf, err := ttmConfig.Load()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	config = conf
}
