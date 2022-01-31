package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type BinStore struct {
	Bins    map[string][]string
	randGen *rand.Rand
}

func NewBinStore() *BinStore {
	source := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(source)

	return &BinStore{make(map[string][]string), gen}
}

func (store *BinStore) NewBin() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[store.randGen.Intn(len(letters))]
	}

	result := string(b)
	store.Bins[result] = []string{}

	return result
}

func (store BinStore) GetBin(binName string) ([]string, bool) {
	bin, exists := store.Bins[binName]
	return bin, exists
}

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

var binStore = NewBinStore()

func main() {
	http.HandleFunc("/r/", binHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome to Rekwest Bin</h1><form method='POST' action='/r/'><button type='submit'>Create a bin</button></form>")
}

func binHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		binName := binStore.NewBin()

		http.Redirect(w, r, "/r/"+binName+"?inspect", 302)
		return
	}

	hash := r.URL.Path[len("/r/"):]
	binAddress := fmt.Sprintf("http://%s/r/%s", r.Host, hash)

	if r.URL.RawQuery == "inspect" {
		rekwests, exists := loadRequest(hash)

		if !exists {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, "<h1>Here are your rekwests:</h1>")
		if len(rekwests) == 0 {
			fmt.Fprintf(w, "<h2>No rekwests</h2>")
			fmt.Fprintf(w, "<p>Make a request to: %s", binAddress)
			fmt.Fprintf(w, "<p>View request at: %s", binAddress+"?inspect")
		}

		for i, rekwest := range rekwests {
			fmt.Fprintf(w, "<h2>Rekwest %d</h2><p>%s</p>", i+1, rekwest)
		}
	} else {
		dump, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		dump = append([]byte("<pre>"), dump...)
		dump = append([]byte(dump), []byte("</pre>")...)

		requesterIP := r.Header.Get("X-Forwarded-For")

		if saveRequest(hash, dump) {
			fmt.Fprintf(w, "<h1>Request saved</h1><p>%s</p>", requesterIP)
			fmt.Fprintf(w, "<p><a href=%s>View requests</a>", binAddress+"?inspect")
		} else {
			http.NotFound(w, r)
		}
	}
}

func saveRequest(hash string, rekwest []byte) bool {
	success := binStore.AddRekwest(hash, string(rekwest))
	return success
}

func loadRequest(hash string) ([]string, bool) {
	bins, exists := binStore.GetBin(hash)
	return bins, exists
}
