package stats

import (
	"io/ioutil"
	"os"

	"github.com/wallix/awless/cloud/mocks"
	"github.com/wallix/awless/database"
)

func init() {
	mocks.InitServices()
}

func newTestDb() (*database.DB, func()) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		panic(e)
	}

	db, err := database.Open(f.Name())
	if err != nil {
		panic(err)
	}

	todefer := func() {
		os.Remove(f.Name())
		db.Close()
	}

	return db, todefer
}
