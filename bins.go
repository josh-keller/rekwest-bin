package main

import (
	"math/rand"
	"time"
)

// Possible letters for the random ID
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Right now the BinStore encapsulates a randome number generator, for...reasons
// Not sure if this is a good idea or not, but it's how it's working for now
type BinStore struct {
	Bins    map[string][]string
	randGen *rand.Rand
}

type RequestInfo struct {
	Raw string
}

type BinInfo struct {
	BinAddress string
	Requests   []RequestInfo
}

// NewBinStore creates a new BinStore and returns a reference to it
// It also seeds the random number generator
func NewBinStore() *BinStore {
	source := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(source)

	return &BinStore{make(map[string][]string), gen}
}

// NewBin creates a new empty bin in the store
func (store *BinStore) NewBin() string {
	// Generate random id: could be factored out
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[store.randGen.Intn(len(letters))]
	}

	result := string(b)
	store.Bins[result] = []string{}

	return result
}

// GetBin gets the bin with the id of binName. It returns an array of
// strings (the requests) and a bool indicating whether the bin exists
func (store BinStore) GetBin(binName string) ([]string, bool) {
	bin, exists := store.Bins[binName]
	return bin, exists
}

// AddRekwest adds a request to the given bin. It returns false if the bin
// does not exist
func (store *BinStore) AddRekwest(binName string, rekwest string) bool {
	_, exists := store.Bins[binName]
	if !exists {
		return false
	}

	binSize := len(store.Bins[binName])

	if binSize >= 20 {
		store.Bins[binName] = store.Bins[binName][binSize+1-20:]
	}

	store.Bins[binName] = append(store.Bins[binName], rekwest)
	return true
}

// Helpers - these could be extracted out along with the BinStore
func saveRequest(hash string, rekwest []byte) bool {
	success := binStore.AddRekwest(hash, string(rekwest))
	return success
}

func loadRequest(hash string) ([]string, bool) {
	bins, exists := binStore.GetBin(hash)
	return bins, exists
}
