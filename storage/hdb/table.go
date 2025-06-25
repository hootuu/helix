package hdb

import "gorm.io/gorm"

type Table struct {
	table string
	model interface{}
}

func NewTable(table string, model interface{}) *Table {
	return &Table{table: table, model: model}
}

func AutoMigrateWithTable(db *gorm.DB, tables ...*Table) error {
	for _, table := range tables {
		err := db.Table(table.table).AutoMigrate(table.model)
		if err != nil {
			return err
		}
	}
	return nil
}
