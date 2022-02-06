package db_controller

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
	client   *mongo.Client
	options  *options.ClientOptions
	bins     *mongo.Collection
	dbName   string
	collName string
}

func NewDatabase(dbName, collName string) *Database {
	err := godotenv.Load()
	checkAndFail(err)

	options := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	return &Database{
		client:   nil,
		options:  options,
		bins:     nil,
		dbName:   dbName,
		collName: collName,
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

const SLICE_SIZE = 3

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
	// Method     string
	// Host       string
	// Path       string
	// Created    string // timestamp
	// Parameters map[string]string
	// Headers    map[string][]string
	// Body       string
	Raw string
}

func NewRekwest(r *http.Request) (Rekwest, error) {
	dump, err := httputil.DumpRequest(r, true)

	checkAndFail(err)

	return Rekwest{
		RekwestId: primitive.NewObjectIDFromTimestamp(time.Now()),
		Raw:       string(dump),
	}, nil

}

func checkAndFail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewBin() (Bin, string) {
	newBin := Bin{
		ObjectID: primitive.NewObjectIDFromTimestamp(time.Now()),
		BinId:    makeRandomId(),
		Rekwests: make([]Rekwest, 0),
	}
	_, err := bins.InsertOne(context.TODO(), newBin)
	checkAndFail(err)

	return newBin, "ok"
}

func FindBin(binId string) (Bin, bool) {
	filter := bson.D{{"binid", binId}}
	var bin Bin
	err := bins.FindOne(context.TODO(), filter).Decode(&bin)

	if err != nil {
		fmt.Println("error: ", err)
		return Bin{}, false
	}

	return bin, true
}

func GetAllBins() {
	var results []Bin
	findOptions := options.Find()

	cursor, err := bins.Find(context.TODO(), bson.D{{}}, findOptions)
	checkAndFail(err)

	for cursor.Next(context.TODO()) {
		var elem Bin

		err := cursor.Decode(&elem)
		checkAndFail(err)
		results = append(results, elem)
	}

	checkAndFail(cursor.Err())

	fmt.Println(results)
	cursor.Close(context.TODO())
}

func AddRekwest(binId string, r *http.Request) error {
	rekwest, err := NewRekwest(r)

	checkAndFail(err)

	result, err := bins.UpdateOne(
		context.TODO(),
		bson.M{"binid": binId},
		bson.M{"$push": bson.M{"rekwests": bson.M{"$each": []Rekwest{rekwest}, "$position": 0, "$slice": SLICE_SIZE}}},
	)

	checkAndFail(err)

	if result.MatchedCount == 0 {
		return errors.New("Bin not found")
	}

	return nil
}
