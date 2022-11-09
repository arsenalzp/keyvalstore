// Hash table storage implementation.

package ht

import (
	"context"
	"gokeyval/internal/server/storage/entity"
	"sync"
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
	i := hash(k)

	var n *Node = ht.table[i]
	if n == nil {
		ht.table[i] = &Node{
			key:  k,
			val:  v,
			next: nil,
		}
		ht.table[i].Lock()
		defer ht.table[i].Unlock()
		ht.size++
		return true, nil
	}

	var prev *Node
	for n != nil {
		if n.key == k {
			n.Lock()
			n.val = v
			defer n.Unlock()
			return true, nil
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
	ht.size++

	return true, nil
}

func (ht *hashTable) Delete(ctx context.Context, k string) (bool, error) {
	i := hash(k)

	n := ht.table[i]
	if n == nil {
		return false, nil
	}

	var prev *Node
	for n != nil {
		if n.key == k {
			if prev == nil {
				ht.table[i].Lock()
				defer ht.table[i].Unlock()
				ht.table[i] = n.next
				ht.size--
				return true, nil
			}

			prev.Lock()
			defer prev.Unlock()
			prev.next = n.next
			ht.size--

			return true, nil
		}
		prev = n
		n = n.next
	}

	return false, nil
}

func (ht *hashTable) Search(ctx context.Context, k string) (string, error) {
	i := hash(k)

	var n *Node = ht.table[i]

	if n == nil {
		return "", nil
	}

	for n != nil {
		if n.key == k {
			return n.val, nil
		}
		n = n.next
	}

	return "", nil
}

func (ht *hashTable) Import(ctx context.Context, data *[]entity.ImportData) (bool, error) {
	for _, i := range *data {
		ht.Insert(ctx, i.Key, i.Value)
	}

	return true, nil
}

func (ht *hashTable) Export(ctx context.Context) (*[]entity.ExportData, error) {
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

	return &exportItems, nil
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
