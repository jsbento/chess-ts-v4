package types

import (
	c "github.com/jsbento/chess-server-v4/internal/engine/constants"
)

type HashEntry struct {
	PosKey uint64
	Move   int
	Score  int
	Depth  int
	Flags  int
}

type HashTable struct {
	HashEntries []HashEntry
	NumEntries  int
	NewWrite    int
	OverWrite   int
	Cut         int
	Hit         int
}

func NewHashTable(size int) (hashTable *HashTable) {
	numEntries := size - 2

	hashTable = &HashTable{
		HashEntries: make([]HashEntry, size),
		NumEntries:  numEntries,
		NewWrite:    0,
		OverWrite:   0,
		Cut:         0,
		Hit:         0,
	}
	hashTable.Clear()

	return
}

func (hashTable *HashTable) Clear() {
	for i := 0; i < hashTable.NumEntries; i++ {
		hashTable.HashEntries[i].PosKey = uint64(0)
		hashTable.HashEntries[i].Move = c.NOMOVE
		hashTable.HashEntries[i].Score = 0
		hashTable.HashEntries[i].Depth = 0
		hashTable.HashEntries[i].Flags = 0
	}

	hashTable.NewWrite = 0
}
