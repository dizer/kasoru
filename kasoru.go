package kasoru

import (
	"fmt"

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
	Cursor uint64
	Limit  uint64
}

// New Kasoru
func New(db *gorm.DB, model interface{}, page Page) (*Kasoru, error) {
	kasoru := Kasoru{
		OriginalDB: db,
		Model:      model,
		Page:       page,
	}

	scope := db.NewScope(model)
	pk := fmt.Sprintf("%s.%s", scope.TableName(), scope.PrimaryKey())

	kasoru.DB = db.
		Where(fmt.Sprintf("%s > ?", pk), page.Cursor).
		Limit(page.Limit).
		Order(fmt.Sprintf("%s ASC", pk))

	return &kasoru, nil
}

// Next Kasoru
func (kasoru *Kasoru) Next(cursor uint64) *Kasoru {
	nextKasoru, _ := New(kasoru.OriginalDB, kasoru.Model, Page{Cursor: cursor, Limit: kasoru.Page.Limit})
	return nextKasoru
}
