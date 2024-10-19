// Define storage interface implementation.
// Initialize underlaying storage.

package storage

import (
	"context"

	"github.com/arsenalzp/keyvalstore/internal/server/errors"
	entity "github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
	ht "github.com/arsenalzp/keyvalstore/internal/server/storage/hash-table"
	sqlite "github.com/arsenalzp/keyvalstore/internal/server/storage/sqlite"
)

// Interface of underlying storage
type Storage interface {
	Search(context.Context, string) (string, error)
	Insert(context.Context, string, string) (bool, error)
	Delete(context.Context, string) (bool, error)
	Import(context.Context, []entity.ImportData) (bool, error)
	Export(context.Context) ([]entity.ExportData, error)
}

// Initialize the underlying storage defined by storage variable
// Returns initialized storage
func NewStrg(kind string) (Storage, error) {
	switch kind {
	case "hash":
		ht, err := ht.NewHT()
		if err != nil {
			return nil, err
		}

		return ht, nil
	case "sqlite":
		db, err := sqlite.NewDb()
		if err != nil {
			return nil, err
		}

		return db, nil
	case "":
		return nil, errors.New("storage type is undefined", errors.StorageKindUndef, nil)
	default:
		return nil, errors.New("storage type is unknown", errors.StorageKindErr, nil)
	}
}
