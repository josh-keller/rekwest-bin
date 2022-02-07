package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client    *mongo.Client
	options   *options.ClientOptions
	bins      *mongo.Collection
	dbName    string
	collName  string
	sliceSize int
}

func NewDatabase(dbName, collName string) *Database {
	err := godotenv.Load()
	checkAndFail(err)
	var opts *options.ClientOptions
	var sliceSize int

	if os.Getenv("TEST_ENV") == "true" {
		opts = options.Client().ApplyURI(os.Getenv("MONGODB_TEST_URI"))
		sliceSize = 3
	} else {
		opts = options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
		sliceSize = 20
	}

	return &Database{
		client:    nil,
		options:   opts,
		bins:      nil,
		dbName:    dbName,
		collName:  collName,
		sliceSize: sliceSize,
	}
}

func (db *Database) Connect() {
	client, err := mongo.Connect(context.TODO(), db.options)
	checkAndFail(err)
	err = client.Ping(context.TODO(), nil)
	checkAndFail(err)

	db.client = client
	db.bins = client.Database(db.dbName).Collection(db.collName)

	fmt.Println("MongoDB connected")
}

func (db *Database) Disconnect() {
	if db.client == nil {
		return
	}

	err := db.client.Disconnect(context.TODO())
	checkAndFail(err)

	fmt.Println("MongoDB disconnected")
}

var makeRandomId = func() func() string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var source = rand.NewSource(time.Now().UnixNano())
	var gen = rand.New(source)

	return func() string {
		b := make([]rune, 8)

		for i := range b {
			b[i] = letters[gen.Intn(len(letters))]
		}

		return string(b)
	}
}()

type Bin struct {
	ObjectID primitive.ObjectID
	BinId    string
	Rekwests []Rekwest
}

func (b Bin) Timestamp() time.Time {
	return b.ObjectID.Timestamp()
}

type Rekwest struct {
	RekwestId primitive.ObjectID
	Method    string
	Host      string
	Path      string
	Params    map[string][]string
	// Headers    map[string][]string
	// Body       string
	Raw string
}

func (rekwest Rekwest) Timestamp() time.Time {
	return rekwest.RekwestId.Timestamp()
}

func NewRekwest(r *http.Request) (Rekwest, error) {
	dump, err := httputil.DumpRequest(r, true)
	checkAndFail(err)
	queryParams := r.URL.Query()

	return Rekwest{
		RekwestId: primitive.NewObjectIDFromTimestamp(time.Now()),
		Method:    r.Method,
		Host:      r.Host,
		Path:      r.URL.Path,
		Raw:       string(dump),
		Params:    queryParams,
	}, nil
}

func checkAndFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (db *Database) NewBin() (Bin, string) {
	newBin := Bin{
		ObjectID: primitive.NewObjectIDFromTimestamp(time.Now()),
		BinId:    makeRandomId(),
		Rekwests: make([]Rekwest, 0),
	}
	_, err := db.bins.InsertOne(context.TODO(), newBin)
	checkAndFail(err)

	return newBin, "ok"
}

func (db *Database) FindBin(binId string) (Bin, bool) {
	filter := bson.D{{"binid", binId}}
	var bin Bin
	err := db.bins.FindOne(context.TODO(), filter).Decode(&bin)

	if err != nil {
		fmt.Println("error: ", err)
		return Bin{}, false
	}

	return bin, true
}

func (db *Database) AddRekwest(binId string, r *http.Request) error {
	rekwest, err := NewRekwest(r)

	checkAndFail(err)

	result, err := db.bins.UpdateOne(
		context.TODO(),
		bson.M{"binid": binId},
		bson.M{"$push": bson.M{"rekwests": bson.M{"$each": []Rekwest{rekwest}, "$position": 0, "$slice": db.sliceSize}}},
	)

	checkAndFail(err)

	if result.MatchedCount == 0 {
		return errors.New("Bin not found")
	}

	return nil
}
