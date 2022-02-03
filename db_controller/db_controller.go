package db_controller

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var bins *mongo.Collection

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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

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

	fmt.Println(results)
	cursor.Close(context.TODO())
}

func AddRekwest(binId string, rekwest Rekwest) bool {
	result, err := bins.UpdateOne(
		context.TODO(),
		bson.M{"binid": binId},
		bson.M{"$push": bson.M{"rekwests": bson.M{"$each": []Rekwest{rekwest}, "$position": 0, "$slice": SLICE_SIZE}}},
	)

	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		return false
	}

	return true
}
