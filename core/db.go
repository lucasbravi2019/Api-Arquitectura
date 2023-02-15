package core

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseConnection *mongo.Database

func GetDatabaseConnection() *mongo.Database {
	if databaseConnection == nil {
		conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

		if err != nil {
			log.Fatal("La conexion a la base de datos no pudo realizarse")
		}

		err = conn.Ping(context.TODO(), nil)

		if err != nil {
			log.Fatal("La conexion a la base de datos no responde")
		}

		databaseConnection = conn.Database("pasteleria")
	}

	return databaseConnection
}

func CheckDatabaseHealth() error {
	log.Println("Checking database health")
	return GetDatabaseConnection().Client().Ping(context.TODO(), nil)
}
