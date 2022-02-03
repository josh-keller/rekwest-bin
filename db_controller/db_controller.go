package db_controller

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var bins *mongo.Collection

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

// clean up structs
type Bin struct {
	BinId      string
	Created_at string // timestamp
	Rekwests   []Rekwest
}

type Rekwest struct {
	// RekwestId  string
	// Method     string
	// Host       string
	// Path       string
	// Created    string // timestamp
	// Parameters map[string]string
	// Headers    map[string]string
	// Body       string
	Raw string
}

func Connect() {
	if client != nil {
		return
	}

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
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
		BinId:      makeRandomId(),
		Created_at: time.Now().GoString(),
		Rekwests:   make([]Rekwest, 0),
	}
	_, err := bins.InsertOne(context.TODO(), newBin)
	if err != nil {
		log.Fatal(err)
	}

	return newBin, "ok"
}

func FindBin(binId string) (Bin, bool) {
	fmt.Println("Finding: ", binId)
	filter := bson.D{primitive.E{Key: "binid", Value: binId}}
	var bin Bin
	err := bins.FindOne(context.TODO(), filter).Decode(&bin)
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

		fmt.Println("error: ", err)
		return Bin{}, false
	}

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

	fmt.Println("Adding: ", bin)

	// objectID, err := primitive.ObjectIDFromHex(binId)
	// if err != nil {
	// 	return false
	// }
	//
	// add slicing functionality
	bin.Rekwests = append(bin.Rekwests, rekwest)

	fmt.Println("Bin.Rekwests: ", bin.Rekwests)

	filter := bson.D{{"binid", binId}}
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
