package kasoru

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

// Kasoru cursor pagination for Gorm
type Kasoru struct {
	OriginalDB *gorm.DB
	DB         *gorm.DB
	Model      interface{}
	Page       Page
}

// Page of pagination
type Page struct {
	Cursor    uint64
	Limit     uint64
	Direction string
}

// New Kasoru
func New(db *gorm.DB, model interface{}, page Page) (*Kasoru, error) {
	kasoru := Kasoru{
		OriginalDB: db,
		Model:      model,
		Page:       page,
	}

	field := kasoru.CursorFieldname()

	direction := "ASC"
	sortingSymbol := ">"
	if page.Direction == "DESC" {
		direction = "DESC"
		sortingSymbol = "<"
	}

	kasoru.DB = db.
		Limit(page.Limit).
		Order(fmt.Sprintf("%s %s", field, direction))

	if page.Cursor != 0 {
		kasoru.DB = kasoru.DB.Where(fmt.Sprintf("%s %s ?", field, sortingSymbol), page.Cursor)
	}

	return &kasoru, nil
}

// Next Kasoru
func (kasoru *Kasoru) Next(cursor uint64) *Kasoru {
	nextKasoru, _ := New(kasoru.OriginalDB, kasoru.Model, Page{Cursor: cursor, Limit: kasoru.Page.Limit, Direction: kasoru.Page.Direction})
	return nextKasoru
}

// CursorFieldname
func (kasoru *Kasoru) CursorFieldname() string {
	tagName := "kasoru"
	scope := kasoru.OriginalDB.NewScope(kasoru.Model)

	s := reflect.Indirect(reflect.ValueOf(kasoru.Model)).Type()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			return fmt.Sprintf("%s.%s", scope.TableName(), gorm.ToDBName(field.Name))
		}
	}

	return fmt.Sprintf("%s.%s", scope.TableName(), scope.PrimaryKey())
}
