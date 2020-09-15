package data

import (
	"github.com/hashicorp/go-memdb"
	"github.com/magazin/data/models"
	"github.com/magazin/data/repository"
)

var db *memdb.MemDB

func Init() error  {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"User": {
				Name: "User",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
						},
				},
			},
			"Product": {
				Name: "Product",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
			"Quantity": {
				Name: "Quantity",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ProductID"},
					},
				},
			},
		},
	}
	newDB, err := memdb.NewMemDB(schema)
	if err != nil {
		return err
	}
	db = newDB

	return insertMockData(db)
}

func insertMockData(db *memdb.MemDB) error {
	users := []*models.User{
		{"987", "123456789", "foo@gmail.com"},
		{"988", "111111111", "bar@address.com"},
		{"999", "222222222", "moo@something.org"},
	}
	txn := db.Txn(true)
	for _, user := range users {
		if err := txn.Insert("User", user); err != nil {
			return err
		}
	}

	products := []*models.Product{
		{"prodID1", "Product11", 112.12 },
		{"prodID2", "Product22", 22.09},
	}
	for _, product := range products {
		if err := txn.Insert("Product", product); err != nil {
			return err
		}
	}

	quants := []*models.Quantity{
		{ProductID: "prodID1", Quantity: 100},
		{ProductID: "prodID2", Quantity: 10},
	}
	for _, quant := range quants {
		if err := txn.Insert("Quantity", quant); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func Begin(write bool) *repository.Transact {
	return repository.NewTransact(db, write)
}
