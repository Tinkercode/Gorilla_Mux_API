package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"gopractice/gorillamux/proto/assignment/assignmentpb"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	db "gopractice/gorillamux/db"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
)

type user struct {
	DocID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID        int                `json:"id,omitempty" bson:"id,omitempty"`
	FirstName string             `json:"Firstname,omitempty" bson:"Firstname,omitempty"`
	LastName  string             `json:"Lastname,omitempty" bson:"Lastname,omitempty"`
	Email     string             `json:"Email,omitempty" bson:"Email,omitempty"`
}

type employee struct {
	DocID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID          int                `json:"id,omitempty" bson:"id,omitempty"`
	UserID      int                `json:"userId,omitempty" bson:"userId,omitempty"`
	Designation string             `json:"designation,omitempty" bson:"designation,omitempty"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/assignment/user", GetHandler).Methods(http.MethodGet)
	r.HandleFunc("/assignment/user", PostHandler).Methods(http.MethodPost)
	r.HandleFunc("/assignment/user", PatchHandler).Methods(http.MethodPatch)
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPatch,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	server := &http.Server{
		Handler:      corsOpts.Handler(r),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	fmt.Println("Starting the API server on ", "http://"+server.Addr, " ........")
	log.Fatal(server.ListenAndServe())
}

//GetHandler for the URL
func GetHandler(resp http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["proto_body"]
	if !ok || len(keys) < 1 {
		log.Fatal("Error in keys")
	}
	key := keys[0]

	bsedat, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		fmt.Println("Error:", err)
	}
	request := &assignmentpb.GetRequest{}
	proto.Unmarshal(bsedat, request)

	dbUser, dbEmployee, client, ctx := db.GetClient()
	var employeeData employee
	var userData user

	userID := request.GetUserId()
	filterEmployee := bson.M{"userId": userID}
	findEmployeeErr := dbEmployee.Collection.FindOne(context.TODO(), filterEmployee).Decode(&employeeData)
	if findEmployeeErr != nil {
		fmt.Println("Error:", findEmployeeErr)
		defer client.Disconnect(ctx)
	} else {
		filterUser := bson.M{"id": employeeData.ID}
		findUserErr := dbUser.Collection.FindOne(context.TODO(), filterUser).Decode(&userData)
		if findUserErr != nil {
			fmt.Println("Error:", findUserErr)
			defer client.Disconnect(ctx)
		} else {
			result := &assignmentpb.GetResponse{ID: int32(userData.ID), Firstname: userData.FirstName, Lastname: userData.LastName, Email: userData.Email, Designation: employeeData.Designation}
			response, err := proto.Marshal(result)
			if err != nil {
				log.Fatalf("Unable to marshal response : %v", err)
			}
			resp.Write(response)
			defer client.Disconnect(ctx)
		}
	}
}

//PostHandler for the URL
func PostHandler(resp http.ResponseWriter, req *http.Request) {
	dbUser, dbEmployee, client, ctx := db.GetClient()
	contentLength := req.ContentLength
	fmt.Printf("Content Length Received : %v\n", contentLength)
	request := &assignmentpb.PostRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Unable to read message from request : %v", err)
		defer client.Disconnect(ctx)
	}
	proto.Unmarshal(data, request)
	fmt.Println("request", request)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	ID := r.Intn(99999-10000) + 10000
	userID := r.Intn(999999-100000) + 100000

	var userData user = user{ID: ID, FirstName: request.GetFirstname(), LastName: request.GetLastname(), Email: request.GetEmail()}
	_, insertUserErr := dbUser.Collection.InsertOne(context.TODO(), userData)
	if insertUserErr != nil {
		fmt.Println("Error:", insertUserErr)
		defer client.Disconnect(ctx)
	} else {
		var employeeData employee = employee{ID: ID, Designation: request.GetDesignation(), UserID: userID}
		_, insertEmployeeErr := dbEmployee.Collection.InsertOne(context.TODO(), employeeData)
		if insertEmployeeErr != nil {
			fmt.Println("Error:", insertEmployeeErr)
			defer client.Disconnect(ctx)
		} else {
			result := &assignmentpb.PostResponse{ID: int32(userData.ID)}
			response, err := proto.Marshal(result)
			if err != nil {
				log.Fatalf("Unable to marshal response : %v", err)
				defer client.Disconnect(ctx)
			}
			resp.Write(response)
		}
	}
}

//PatchHandler for the URL
func PatchHandler(resp http.ResponseWriter, req *http.Request) {
	dbUser, dbEmployee, client, ctx := db.GetClient()
	contentLength := req.ContentLength
	fmt.Printf("Content Length Received : %v\n", contentLength)

	keys, ok := req.URL.Query()["proto_body"]
	if !ok || len(keys) < 1 {
		log.Fatal("Error in keys")
	}
	key := keys[0]

	bsedat, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		fmt.Println("Error:", err)
	}
	request := &assignmentpb.PatchRequest{}
	proto.Unmarshal(bsedat, request)

	var employeeData employee
	filterEmployee := bson.M{"userId": request.GetUserId()}
	findEmployeeErr := dbEmployee.Collection.FindOne(context.TODO(), filterEmployee).Decode(&employeeData)
	if findEmployeeErr != nil {
		fmt.Println("Error:", findEmployeeErr)
		defer client.Disconnect(ctx)
	} else {
		var userData user = user{Email: request.GetEmail()}
		filterUser := bson.M{"id": bson.M{"$eq": employeeData.ID}}
		update := bson.M{"$set": &userData}
		_, updateEmployeeErr := dbUser.Collection.UpdateOne(context.TODO(), filterUser, update)
		if updateEmployeeErr != nil {
			fmt.Println("Error:", findEmployeeErr)
			defer client.Disconnect(ctx)
		}
	}
}
