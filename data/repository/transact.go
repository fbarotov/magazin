package repository

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"github.com/magazin/data/models"
)

var ErrNotFound = errors.New("not found")

func NewTransact(db *memdb.MemDB, write bool) *Transact {
	return &Transact{db.Txn(write)}
}

type Transact struct {
	txn *memdb.Txn
}

func (t *Transact) Commit()  {
	t.txn.Commit()
}

func (t *Transact) Abort() {
	t.txn.Abort()
}

func (t *Transact) User(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user id cannot be null")
	}

	user, err :=  t.txn.First("User", "id", id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user.(*models.User), nil
}

func (t *Transact) Product(id string)(*models.Product, error) {
	product, err := t.txn.First("Product", "id", id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrNotFound
	}
	return product.(*models.Product), nil
}

func (t *Transact) Quantity(prodID string) (*models.Quantity, error) {
	result, err := t.txn.First("Quantity", "id", prodID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrNotFound
	}
	return result.(*models.Quantity), nil
}

func (t *Transact) InsertQuantity(quantity models.Quantity) error {
	return t.txn.Insert("Quantity", &quantity)
}