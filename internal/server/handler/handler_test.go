package handler

import (
	"context"
	cli "gokeyval/internal/cli/command"
	entity "gokeyval/internal/server/storage/entity"
	"net"
	"reflect"
	"testing"
)

const KEY = "key100000"
const VALUE = "value100000"

type Storage struct {
	storage map[string]string
}

func initStorage() *Storage {
	stg := &Storage{}

	stg.storage = make(map[string]string)

	// for i := 0; i < 100; i++ {
	// 	key := "key" + fmt.Sprint(i)
	// 	value := "valu" + fmt.Sprint(i)

	// 	stg.storage[key] = value
	// }
	return stg
}

func TestDelHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	stg.storage = map[string]string{KEY: VALUE}

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	err := cli.Del(clientConn, nil, []string{KEY})
	if err != nil {
		t.Errorf("error in Del command: %s", err)
		return
	}

	if value, ok := stg.storage[KEY]; ok {
		t.Errorf("error deleting the key, expected: %s, got: %s\n", "\"\"", value)
		return
	}
}
func TestSetHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	err := cli.Set(clientConn, nil, []string{KEY + "=" + VALUE})
	if err != nil {
		t.Errorf("error in Set command: %s", err)
		return
	}

	if value := stg.storage[KEY]; value != VALUE {
		t.Errorf("error getting value in Set Handler, expected: %s, got: %s\n", VALUE, value)
		return
	}
}

func TestGetHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	stg.storage = map[string]string{KEY: VALUE}

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	data, err := cli.Get(clientConn, nil, []string{KEY})
	if err != nil {
		t.Errorf("error in Get command: %s", err)
		return
	}

	if !reflect.DeepEqual(data, []byte(VALUE)) {
		t.Errorf("error getting value in Get Handler: expected: %s got: %s\n", VALUE, data)
		return
	}
}

func (s *Storage) Search(ctx context.Context, key string) (string, error) {
	return s.storage[key], nil
}

func (s *Storage) Insert(ctx context.Context, key string, value string) (bool, error) {
	s.storage[key] = value
	if v, ok := s.storage[key]; ok && v == value {
		return true, nil
	}

	return false, nil
}

func (s *Storage) Delete(ctx context.Context, key string) (bool, error) {
	delete(s.storage, key)

	return true, nil
}

func (s *Storage) Import(context.Context, []entity.ImportData) (bool, error) {
	return true, nil
}
func (s *Storage) Export(context.Context) ([]entity.ExportData, error) {
	return nil, nil
}
