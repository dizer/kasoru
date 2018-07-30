# Kasoru - cursor pagination lib for Go and Gorm

```go
import (
  "github.com/dizer/kasoru"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

db, _ := gorm.Open("sqlite3", "tmp/kasoru.db")
kasoru, _ := New(db, &GormModel{}, kasoru.Page{Cursor: 0, Limit: 2})

// Now you can use kasoru.DB
kasoru.DB.Where("something IS NOT NULL").Find(&something)

// Next page
var ids []uint
kasoru.DB.Model(&GormModel{}).Pluck("id", &ids)
kasoru = kasoru.Next(uint64(ids[len(ids)-1]))
kasoru.Page // {Cursor: 5, Limit: 2}
```

## Using not primary key for cursor

Just use `kasoru` tag:

```go
type Model struct {
	ID       uint `gorm:"primary_key"`
	Position uint `kasoru:"cursor"`
}
```
