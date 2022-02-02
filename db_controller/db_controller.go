package db_controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var bins *mongo.Collection

// clean up structs
type Bin struct {
	BinId      string
	Created_at string // timestamp
	Rekwests   []Rekwest
}

type Rekwest struct {
	RekwestId  string
	Method     string
	Host       string
	Path       string
	Created    string // timestamp
	Parameters map[string]string
	Headers    map[string]string
	Body       string
	Raw        string
}

func Connect() {
	if client != nil {
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_TEST_URI"))

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected")
	bins = client.Database("rekwest-bin").Collection("bins")
}

func Disconnect() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB disconnected")
}

// update how we create the binId / url hash
func NewBin() (Bin, string) {
	newBin := Bin{
		BinId:      "",
		Created_at: time.Now().GoString(),
		Rekwests:   make([]Rekwest, 0),
	}
	results, err := bins.InsertOne(context.TODO(), newBin)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(newBin)
	newBin.BinId = results.InsertedID.(primitive.ObjectID).Hex()
	// fmt.Println(newBin)
	return newBin, newBin.BinId
}

func FindBin(binId string) (Bin, bool) {
	objectID, err := primitive.ObjectIDFromHex(binId)
	if err != nil {
		return Bin{}, false
	}
	filter := bson.D{{"_id", objectID}}
	var bin Bin
	err = bins.FindOne(context.TODO(), filter).Decode(&bin)
	if err != nil {
		// // check if these err checks are needed
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		// if err == mongo.ErrNoDocuments {
		// 	return Bin{}, true
		// } else {
		//   log.Fatal(err)
		//   return Bin{}, false
		// }

		return Bin{}, false
	}

	fmt.Println(bin)
	bin.BinId = binId
	fmt.Println(bin)

	return bin, true
}

func GetAllBins() {
	var results []Bin
	findOptions := options.Find()

	cursor, err := bins.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var elem Bin
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// not returning anything, just printing
	fmt.Println(results)
	cursor.Close(context.TODO())
}

func AddRekwest(binId string, rekwest Rekwest) bool {
	bin, success := FindBin(binId)
	if !success {
		return false
	}

	fmt.Println(bin)

	objectID, err := primitive.ObjectIDFromHex(binId)
	if err != nil {
		return false
	}

	// add slicing functionality
	bin.Rekwests = append(bin.Rekwests, rekwest)

	filter := bson.D{{"_id", objectID}}
	update := bson.D{{"$set", bson.D{{"rekwests", bin.Rekwests}}}}
	result, err := bins.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		return false
	}

	return true
}

/*

	"github.com/wboard82/rekwest-bin/db_controller"

var testBin = db_controller.Bin{
	BinId:      "",
	Created_at: time.Now().GoString(), // timestamp
	Rekwests:   make([]db_controller.Rekwest, 20),
}

var testRekwest = db_controller.Rekwest{
	RekwestId:  "",
	Method:     "POST",
	Host:       "316e-174-81-238-56.ngrok.io",
	Path:       "/r/",
	Created:    time.Now().GoString(), // timestamp
	Parameters: nil,
	Headers: map[string]string{
		"User-Agent":        "curl/7.68.0",
		"Content-Length":    "28",
		"Accept":            "*\/*",
		"Accept-Encoding":   "gzip",
		"Content-Type":      "application/json",
		"X-Forwarded-For":   "192.222.245.48",
		"X-Forwarded-Proto": "https",
	},
	Body: `{"dragons": "are dangerous"}`,
	Raw:  "hi im a raw rekwest",
}

func main() {
	db_controller.Connect()
	defer db_controller.Disconnect()
	bin, binId := db_controller.NewBin()
	fmt.Println(bin, binId, bin.BinId)
	bin, success := db_controller.FindBin(binId)
	fmt.Println(bin, success)
	db_controller.GetAllBins()
	db_controller.AddRekwest(binId, testRekwest)
}
*/
