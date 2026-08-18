// Command migrationhash recalculates the Atlas migration directory checksum.
package main

import (
	"log"

	"ariga.io/atlas/sql/migrate"
)

func main() {
	dir, err := migrate.NewLocalDir("ent/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}
	sum, err := dir.Checksum()
	if err != nil {
		log.Fatal(err)
	}
	if err := migrate.WriteSumFile(dir, sum); err != nil {
		log.Fatal(err)
	}
}
