package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vishal21121/myapp/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "userAuth"
const collectionName = "User"

var collection *mongo.Collection

func Init() {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")
	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("Collection instance created")
}

// function to add the user account in the database
func insertOneUser(user model.User) bool {
	var userGot bson.M
	filter := bson.M{"email": user.Email}
	result := collection.FindOne(context.Background(), filter)
	_ = result.Decode(&userGot)
	if userGot == nil {
		fmt.Println("empty")
		fmt.Println(userGot)
		inserted, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted to the database with the id: ", inserted.InsertedID)
		return true
	}
	return false

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	result := insertOneUser(user)
	if result{
		json.NewEncoder(w).Encode(user)
		return
	}
	json.NewEncoder(w).Encode("A user with this email already exits")
	
}
