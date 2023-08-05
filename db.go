package store

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// MongoDatabase 数据库连接
	MongoDatabase *mongo.Database
)

// New 创建新连接池
func New(uri, dbName string) {
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	MongoDatabase = client.Database(dbName)
}
