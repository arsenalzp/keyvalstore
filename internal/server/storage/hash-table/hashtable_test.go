package ht

import (
	"context"
	"fmt"
	"testing"

	"github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
)

const KEY = "key100000"
const VALUE = "value100000"

func TestInsert(t *testing.T) {
	ctx := context.Background()

	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	// load the test KEY-VALUE pair into the storage
	_, err = hashTbale.Insert(ctx, KEY, VALUE)
	if err != nil {
		t.Errorf("error inserting key-value data: %s\n", err)
		return
	}

	data, err := hashTbale.Search(ctx, KEY)
	if err != nil {
		t.Errorf("error searching a key data: %s\n", err)
		return
	}

	if data != VALUE {
		t.Errorf("error searching a key data: expected %s, got %s\n", VALUE, data)
		return
	}
}

func TestSearch(t *testing.T) {
	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	testSet := populateTestSet(1000)

	// load the test set data into the storage
	for k, v := range testSet {
		_, err := hashTbale.Insert(ctx, k, v)
		if err != nil {
			t.Errorf("error inserting key-value data: %s\n", err)
			return
		}
	}

	for k, v := range testSet {
		result, err := hashTbale.Search(ctx, k)
		if err != nil {
			t.Errorf("error searching for key data: %s\n", err)
			return
		}

		if v != result {
			t.Errorf("error searching for key data: expected %s, got %s\n", v, result)
			return
		}
	}
}

func TestDelete(t *testing.T) {
	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	testSet := populateTestSet(1000)

	// load the test set data into the storage
	for k, v := range testSet {
		_, err := hashTbale.Insert(ctx, k, v)
		if err != nil {
			t.Errorf("error inserting key-value data: %s\n", err)
			return
		}
	}

	// delete loaded data
	for k := range testSet {
		_, err := hashTbale.Delete(ctx, k)
		if err != nil {
			t.Errorf("error deleting key: %s\n", err)
			return
		}
	}

	// test for deletion
	for k := range testSet {
		result, err := hashTbale.Search(ctx, k)
		if err != nil {
			t.Errorf("error searching for key: %s\n", err)
			return
		}
		if result != "" {
			t.Errorf(`error deleting key, expected "", got %s\n`, result)
			return
		}
	}
}

func TestImport(t *testing.T) {
	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	importData := populateImportData(1000)

	_, err = hashTbale.Import(ctx, importData)
	if err != nil {
		t.Errorf("error importing data: %s\n", err)
		return
	}

	// test imported data
	for _, i := range importData {
		result, err := hashTbale.Search(ctx, i.Key)
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
	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	ctx := context.Background()

	// create the test set
	testSet := populateTestSet(1000)

	// load the test set data into the storage
	for k, v := range testSet {
		_, err := hashTbale.Insert(ctx, k, v)
		if err != nil {
			t.Errorf("error inserting key-value data: %s\n", err)
			return
		}
	}

	// export data from the storage
	result, err := hashTbale.Export(ctx)
	if err != nil {
		t.Errorf("error exporting key-value data: %s\n", err)
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

func TestNewHt(t *testing.T) {
	hashTbale, err := NewHT()
	if err != nil {
		t.Errorf("error creating hash table storage: %s\n", err)
		return
	}

	if hashTbale == nil {
		t.Errorf("error creating hash table storage: hash table storage wasn't initialized")
		return
	}
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

func populateImportData(count int) []entity.ImportData {
	var testSet []entity.ImportData

	for i := 1; i <= count; i++ {
		key := "key" + fmt.Sprint(i)
		val := "val" + fmt.Sprint(i)
		testSet = append(testSet, entity.ImportData{key, val})
	}

	return testSet
}
