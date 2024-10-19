package sqlite

import (
	"context"
	"fmt"
	"os"
	"testing"

	entity "github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
)

const KEY = "key100000"
const VALUE = "value100000"

func TestNewDB(t *testing.T) {
	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	if db == nil {
		t.Errorf("error creating DB storage: DB storage wasn't initialized")
		return
	}

	if _, err := os.Stat("default.db"); err != nil {
		t.Errorf("error creating DB storage, %s\n", err)
		return
	}
}

func TestInsert(t *testing.T) {
	defer cleanUp()

	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	ctx := context.Background()

	_, err = db.Insert(ctx, KEY, VALUE)
	if err != nil {
		t.Errorf("error inserting into DB storage: %s\n", err)
		return
	}

	result, err := db.Search(ctx, KEY)
	if err != nil {
		t.Errorf("error searching in DB storage: %s\n", err)
		return
	}

	if result != VALUE {
		t.Errorf("error searching a key in DB storage, expected %s, got %s\n", VALUE, result)
		return
	}
}

func TestDelete(t *testing.T) {
	defer cleanUp()

	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	ctx := context.Background()

	_, err = db.Insert(ctx, KEY, VALUE)
	if err != nil {
		t.Errorf("error inserting into DB storage: %s\n", err)
		return
	}

	result, err := db.Delete(ctx, KEY)
	if err != nil {
		t.Errorf("error deleting from DB storage: %s\n", err)
		return
	}

	if !result {
		t.Errorf("error deleting from DB storage, result expected to be true, got %t\n", result)
		return
	}
}

func TestImport(t *testing.T) {
	defer cleanUp()

	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	importData := populateImportData(1000)

	_, err = db.Import(ctx, importData)
	if err != nil {
		t.Errorf("error importing data: %s\n", err)
		return
	}

	// test imported data
	for _, i := range importData {
		result, err := db.Search(ctx, i.Key)
		if err != nil {
			t.Errorf("error searching for key: %s\n", err)
			return
		}
		if result[3:] != i.Key[3:] {
			t.Errorf("error importing key, expected val%s, got val%s\n", result[3:], i.Key[3:])
			return
		}
	}
}

func TestExport(t *testing.T) {
	defer cleanUp()

	db, err := NewDb()
	if err != nil {
		t.Errorf("error creating DB storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	testSet := populateTestSet(1000)

	// load the test set data into the storage
	for k, v := range testSet {
		_, err := db.Insert(ctx, k, v)
		if err != nil {
			t.Errorf("error inserting key-value pair into DB storage: %s\n", err)
			return
		}
	}

	// export data from the storage
	result, err := db.Export(ctx)
	if err != nil {
		t.Errorf("error exporting key-value data from DB storage: %s\n", err)
		return
	}

	// test exported data with the test set
	for _, i := range result {
		if v := testSet[i.Key]; v != i.Value {
			t.Errorf("error searching for key in the exported data: expected %s, got %s\n", v, i.Value)
			return
		}
	}
}

func cleanUp() error {
	err := os.Remove("default.db")
	if err != nil {
		return err
	}

	return nil
}

func populateImportData(count int) []entity.ImportData {
	var testSet []entity.ImportData

	for i := 1; i <= count; i++ {
		key := "key" + fmt.Sprint(i)
		val := "val" + fmt.Sprint(i)
		testSet = append(testSet, entity.ImportData{key, val})
	}

	return testSet
}

func populateTestSet(count int) map[string]string {
	var testSet map[string]string = make(map[string]string)

	for i := 1; i <= count; i++ {
		key := "key" + fmt.Sprint(i)
		val := "val" + fmt.Sprint(i)
		testSet[key] = val
	}

	return testSet
}
