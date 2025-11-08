package main

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
)

func main() {
	envx.Load()

	db := dbx.New()

	if err := db.AutoMigrate(); err != nil {
		panic(err)
	}
}
