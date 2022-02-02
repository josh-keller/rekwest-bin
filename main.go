package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

var templates = template.Must(template.ParseFiles("templates/inspect.html"))

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
	// The POST route creates a new bin and redirects to the inspect page
	if r.Method == "POST" {
		binName := binStore.NewBin()

		http.Redirect(w, r, "/r/"+binName+"?inspect", 302)
		return
	}

	// This grabs the part after /r/ in the path
	binID := r.URL.Path[len("/r/"):]
	// Put the full link together here to be displayed on a landing page
	binAddress := fmt.Sprintf("http://%s/r/%s", r.Host, binID)

	// If there is a query "inspect", show all the requests
	if r.URL.RawQuery == "inspect" {
		rekwests, exists := loadRequest(binID)

		if !exists {
			http.NotFound(w, r)
			return
		}

		requestInfo := make([]RequestInfo, len(rekwests))

		for i, req := range rekwests {
			requestInfo[i] = RequestInfo{req}
		}

		bin := BinInfo{
			BinAddress: binAddress,
			Requests:   requestInfo,
		}

		renderTemplate(w, "inspect", &bin)

	} else {
		dump, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		requesterIP1 := r.RemoteAddr
		requesterIP := r.Header.Get("X-Forwarded-For")

		if saveRequest(binID, dump) {
			fmt.Fprintf(w, "<h1>Request saved</h1><p>%s, %s</p>", requesterIP1, requesterIP)
			fmt.Fprintf(w, "<p><a href=%s>View requests</a>", binAddress+"?inspect")
		} else {
			http.NotFound(w, r)
		}
	}
}

func renderTemplate(writer http.ResponseWriter, tmpl string, bin *BinInfo) {
	err := templates.ExecuteTemplate(writer, tmpl+".html", bin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
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
