package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"log"
	"testing"
)

func TestGetConfigurationValues(t *testing.T) {
	Convey("Given an arbitrary file...", t, func() {
		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		//defer os.Remove(tmpFile.Name())
		log.Printf("Creating temporary file at: %s", tmpFile.Name())

		Convey("Should be able to create a database there...", func() {
			err := createDatabaseForTesting(tmpFile.Name())
			So(err, ShouldBeNil)

			Convey("Should be able to get all database configuration keys", func() {
				db, err := sqlx.Open("sqlite3", tmpFile.Name())
				So(err, ShouldBeNil)
				defer db.Close()
				keys, err := GetConfigurationValues(db)
				So(err, ShouldBeNil)
				So(len(keys), ShouldEqual, 1)

				So(keys[0].Key, ShouldEqual, "db_schema")
				So(keys[0].Value, ShouldEqual, "v1")
			})
		})
	})
}

func TestGetConfigurationValue(t *testing.T) {
	Convey("Given an arbitrary file...", t, func() {
		tmpFile, err := ioutil.TempFile("", "repo")
		So(err, ShouldBeNil)
		//defer os.Remove(tmpFile.Name())
		log.Printf("Creating temporary file at: %s", tmpFile.Name())

		Convey("Should be able to create a database there...", func() {
			err := createDatabaseForTesting(tmpFile.Name())
			So(err, ShouldBeNil)

			Convey("Should be able to get a configuration value...", func() {
				db, err := sqlx.Open("sqlite3", tmpFile.Name())
				So(err, ShouldBeNil)
				defer db.Close()
				key, err := GetConfigurationValue("db_schema", db)
				So(err, ShouldBeNil)

				So(key.Key, ShouldEqual, "db_schema")
				So(key.Value, ShouldEqual, "v1")
			})
		})
	})
}
