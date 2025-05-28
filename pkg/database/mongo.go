package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"time"
)

var MongoDB *DBConnection

// DBConnection defines the connection structure
type DBConnection struct {
	session *mgo.Session
}

// NewConnection handles connecting to a mongo database
func NewConnection(host, dbName, user, pwd string) (conn *DBConnection) {
	info := &mgo.DialInfo{
		// Address if its a local db then the value host=localhost
		Addrs: []string{host},
		// Timeout when a failure to connect to db
		Timeout: 60 * time.Second,
		// Database name
		Database: dbName,
		// Database credentials if your db is protected
		Username: user,
		Password: pwd,
	}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	conn = &DBConnection{session}
	return conn
}

// Use handles connect to a certain collection
func (conn *DBConnection) Use(dbName, tableName string) (collection *mgo.Collection) {
	// This returns method that interacts with a specific collection and table
	return conn.session.DB(dbName).C(tableName)
}

// Close handles closing a database connection
func (conn *DBConnection) Close() {
	// This closes the connection
	conn.session.Close()
	return
}

// var MongoDBClient *mongo.Client
var MongoDB1 MongoDBClient

type MongoDBClient struct {
	Client *mongo.Client
}

func NewConnectionV1(host string) {
	uri := fmt.Sprintf("mongodb://%s", host)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err.Error())
	}
	MongoDB1.Client = client
}

// Find
// @Description: 查询多条数据
func (db *MongoDBClient) Find(database, collection string, selector bson.M, opts *options.FindOptions, schema interface{}) error {
	c := db.Client.Database(database).Collection(collection)
	//var opts = &options.FindOptions{}
	//if pageSize > 0 {
	//	opts.SetLimit(pageSize)
	//	opts.SetSkip((pageSize * page) - pageSize)
	//}
	cur, err := c.Find(context.Background(), selector, opts)
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())
	return cur.All(context.Background(), schema)
}

func (db *MongoDBClient) Count(database, collection string, selector bson.M) int64 {
	c := db.Client.Database(database).Collection(collection)
	//var opts = &options.FindOptions{}
	count, err := c.CountDocuments(context.Background(), selector)
	if err != nil {
		fmt.Println(err.Error())
	}

	return count
}

func (db *MongoDBClient) Sum(database, collection, amountField string, filter bson.M) (int, error) {
	c := db.Client.Database(database).Collection(collection)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$" + amountField}}}},
	}

	cursor, err := c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		TotalAmount int `bson:"total"`
	}

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.TotalAmount, nil
}
