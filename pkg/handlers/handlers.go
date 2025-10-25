package handlers

import (
	ttmStore "ttm/pkg/store"
	ttmDB "ttm/pkg/store/db"
)

var store *ttmStore.Store

func init() {
	store = ttmStore.NewStore(ttmDB.NewDBStore())
	store.Init()
}
