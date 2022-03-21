package main

import (
	"github/com/MonikaCat/ibc-token/database"
	"log"

)

func main() {

	db, err := database.GetDatabase()
	if err != nil {
		log.Fatal("error while getting database: ", err)
	}
	defer db.Sqlx.Close()

}
