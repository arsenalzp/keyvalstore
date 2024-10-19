// Hash table storage implementation.

package ht

import (
	"context"
	"sync"

	"github.com/arsenalzp/keyvalstore/internal/server/errors"
	"github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
)

// Size of hash table
const _HT_SIZE uint32 = 1048573

type Node struct {
	sync.Mutex

	key  string
	val  string
	next *Node
}

// var hashTable []*Node
type hashTable struct {
	table []*Node
	size  uint64
}

func (ht *hashTable) Insert(ctx context.Context, k, v string) (bool, error) {
	dataCh := make(chan struct{}, 1)

	go func(h *hashTable, c chan<- struct{}, k string) {
		i := hash(k)

		var n *Node = ht.table[i]
		if n == nil {
			h.table[i] = &Node{
				key:  k,
				val:  v,
				next: nil,
			}
			h.table[i].Lock()
			defer h.table[i].Unlock()
			h.size++
			c <- struct{}{}
			return
		}

		var prev *Node
		for n != nil {
			if n.key == k {
				n.Lock()
				n.val = v
				defer n.Unlock()
				c <- struct{}{}
				return
			}
			prev = n
			n = n.next
		}

		prev.next = &Node{
			key:  k,
			val:  v,
			next: nil,
		}
		prev.next.Lock()
		defer prev.next.Unlock()
		h.size++

		c <- struct{}{}
	}(ht, dataCh, k)

	select {
	case <-ctx.Done():
		err := errors.New("hash table error: canceled", errors.HashTabInsErr, nil)
		return false, err
	case <-dataCh:
		return true, nil
	}
}

func (ht *hashTable) Delete(ctx context.Context, k string) (bool, error) {
	dataCh := make(chan struct{}, 1)

	go func(h *hashTable, c chan<- struct{}, k string) {
		i := hash(k)

		n := h.table[i]
		if n == nil {
			c <- struct{}{}
			return
		}

		var prev *Node
		for n != nil {
			if n.key == k {
				if prev == nil {
					h.table[i].Lock()
					defer h.table[i].Unlock()
					h.table[i] = n.next
					h.size--
					c <- struct{}{}
					return
				}

				prev.Lock()
				defer prev.Unlock()
				prev.next = n.next
				h.size--
				c <- struct{}{}
				return
			}
			prev = n
			n = n.next
		}
	}(ht, dataCh, k)

	select {
	case <-ctx.Done():
		err := errors.New("hash table error: canceled", errors.HashTabDelErr, nil)
		return false, err
	case <-dataCh:
		return true, nil
	}
}

func (ht *hashTable) Search(ctx context.Context, k string) (string, error) {
	dataCh := make(chan string, 1)

	go func(h *hashTable, c chan<- string, k string) {
		i := hash(k)

		var n *Node = ht.table[i]

		if n == nil {
			c <- ""
		}

		for n != nil {
			if n.key == k {
				c <- n.val
			}
			n = n.next
		}

		c <- ""
	}(ht, dataCh, k)

	select {
	case <-ctx.Done():
		err := errors.New("hash table error: canceled", errors.HashTabSrchErr, nil)
		return "", err
	case value := <-dataCh:
		return value, nil
	}
}

func (ht *hashTable) Import(ctx context.Context, data []entity.ImportData) (bool, error) {
	for _, i := range data {
		_, err := ht.Insert(ctx, i.Key, i.Value)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (ht *hashTable) Export(ctx context.Context) ([]entity.ExportData, error) {
	dataCh := make(chan []entity.ExportData, 1)

	go func(h *hashTable, c chan<- []entity.ExportData) {
		var exportItems []entity.ExportData

		for _, n := range ht.table {
			if n == nil {
				continue
			}
			for n != nil {
				exportItems = append(exportItems, entity.ExportData{n.key, n.val})
				n = n.next
			}
		}

		c <- exportItems
	}(ht, dataCh)

	select {
	case <-ctx.Done():
		err := errors.New("hash table error: canceled", errors.HashTabExpErr, nil)
		return nil, err
	case value := <-dataCh:
		return value, nil
	}
}

// Calculate hash function for a string
func hash(str string) uint32 {
	var hash uint32
	var power uint32 = 1
	for _, char := range str {
		char_code := uint32(char - rune('!') + 1)
		hash = (hash + power*char_code) % _HT_SIZE
		power = (power * 131) % _HT_SIZE
	}
	return hash
}

func NewHT() (*hashTable, error) {
	storage := &hashTable{
		table: make([]*Node, _HT_SIZE),
	}

	return storage, nil
}
