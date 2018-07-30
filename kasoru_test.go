package kasoru

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Model struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

type Model2 struct {
	ID       uint `gorm:"primary_key"`
	Position uint `kasoru:"cursor"`
	Name     string
}

func TestField(t *testing.T) {
	os.Remove("tmp/TestField.db")
	db, err := gorm.Open("sqlite3", "tmp/TestField.db")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	var kasoru *Kasoru

	kasoru, _ = New(db, &Model{}, Page{Cursor: 0, Limit: 2})
	if kasoru.CursorFieldname() != "models.id" {
		t.Errorf("CursorFieldname ignores primary key")
	}

	kasoru, _ = New(db, &Model2{}, Page{Cursor: 0, Limit: 2})
	if kasoru.CursorFieldname() != "model2.cursor" {
		t.Errorf("CursorFieldname ignores tag")
	}
}

func TestNew(t *testing.T) {
	db, err := gorm.Open("sqlite3", "tmp/TestNew.db")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	_, err = New(db, &Model{}, Page{})

	if err != nil {
		t.Error(err)
	}
}

func TestLimit(t *testing.T) {
	os.Remove("tmp/TestLimit.db")
	db, err := gorm.Open("sqlite3", "tmp/TestLimit.db")
	defer db.Close()
	db.AutoMigrate(&Model{})

	db.Create(&Model{Name: "a"})
	db.Create(&Model{Name: "b"})
	db.Create(&Model{Name: "c"})

	if err != nil {
		panic(err)
	}

	kasoru, _ := New(db, &Model{}, Page{Cursor: 0, Limit: 2})

	var kasoruNames []string
	kasoru.DB.Model(&Model{}).Pluck("name", &kasoruNames)

	var names []string
	db.Model(&Model{}).Pluck("name", &names)

	if strings.Join(kasoruNames, ":") != "a:b" {
		t.Errorf("kasoru doesnt limit rows")
	}

	if strings.Join(names, ":") != "a:b:c" {
		t.Errorf("insert failed")
	}
}

func TestNext(t *testing.T) {
	os.Remove("tmp/TestNext.db")
	db, err := gorm.Open("sqlite3", "tmp/TestNext.db")
	defer db.Close()
	db.AutoMigrate(&Model{})

	db.Create(&Model{Name: "a"})
	db.Create(&Model{Name: "b"})
	db.Create(&Model{Name: "c"})

	if err != nil {
		panic(err)
	}

	kasoru, _ := New(db, &Model{}, Page{Cursor: 0, Limit: 2})
	var ids []uint
	kasoru.DB.Model(&Model{}).Pluck("id", &ids)
	nextKasoru := kasoru.Next(uint64(ids[len(ids)-1]))

	if nextKasoru.Page.Cursor != uint64(2) {
		t.Errorf("next cursor should be equal to last of current")
	}

	if nextKasoru.Page.Limit != uint64(2) {
		t.Errorf("next cursor should keep limit")
	}
}

func TestFlow(t *testing.T) {
	os.Remove("tmp/TestFlow.db")
	db, err := gorm.Open("sqlite3", "tmp/TestFlow.db")
	defer db.Close()
	db.AutoMigrate(&Model{})

	db.Create(&Model{Name: "a"})
	db.Create(&Model{Name: "b"})
	db.Create(&Model{Name: "c"})
	db.Create(&Model{Name: "d"})
	db.Create(&Model{Name: "e"})

	if err != nil {
		panic(err)
	}

	kasoru, _ := New(db, &Model{}, Page{Cursor: 0, Limit: 2})
	var ids []uint

	for len(ids) < 5 {
		kasoru.DB.Model(&Model{}).Pluck("id", &ids)
		kasoru = kasoru.Next(uint64(ids[len(ids)-1]))
	}

	if fmt.Sprint(ids) != "[1 2 3 4 5]" {
		t.Errorf("kasoru should iterate over db")
	}
}
