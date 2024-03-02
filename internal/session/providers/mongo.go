package providers

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/eyko139/go-snippets/internal/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pder = &MongoSessionProvider{}

type MongoSessionStore struct {
	sid          string                      `bson:"sid"`
	timeAccessed time.Time                   `bson:"timeAccessed"`
	value        map[interface{}]interface{} `bson:"value"`
}

type MongoSessionProvider struct {
	lock       sync.Mutex
	collection *mongo.Collection
	sessions   map[string]*list.Element
}

func (mss *MongoSessionStore) Set(key, value interface{}) error {
	mss.value[key] = value
	pder.SessionUpdate(mss.sid, mss.value)
	return nil
}

func (mss *MongoSessionStore) Get(key interface{}) interface{} {
	result := pder.collection.FindOne(context.TODO(), key)
	return result
}

func (mss *MongoSessionStore) Delete(key interface{}) error {
	_, err := pder.collection.DeleteOne(context.TODO(), key)
	return err
}

func (mss *MongoSessionStore) SessionID() string {
	return ""
}

func (msp *MongoSessionProvider) SessionInit(sid string) (session.Session, error) {
	v := make(map[interface{}]interface{})
	newSession := &MongoSessionStore{sid: sid, timeAccessed: time.Now().Local(), value: v}
	_, err := msp.collection.InsertOne(context.TODO(), bson.D{{"sid", newSession.sid}, {"timeAccessed", newSession.timeAccessed}, {"value", v}})
	if err != nil {
		panic(err)
	}
	return newSession, nil
}

func (msp *MongoSessionProvider) SessionRead(sid string) (session.Session, error) {
	var result MongoSessionStore
	err := msp.collection.FindOne(context.TODO(), bson.D{{"sid", sid}}).Decode(&result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}

func (msp *MongoSessionProvider) SessionDestroy(sid string) error {
	_, err := msp.collection.DeleteOne(context.TODO(), sid)
	return err
}
func (msp *MongoSessionProvider) SessionGC(maxLifeTime int64) {
	//TODO: implement
}

func (msp *MongoSessionProvider) SessionUpdate(sid string, update interface{}) error {
	_, err := msp.collection.UpdateOne(context.TODO(), bson.D{{"sid", sid}}, update)
	return err
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:password@localhost:27017"))
	if err != nil {
		panic(err)
	}
	collection := client.Database("snippets").Collection("sessions")
	pder.collection = collection
	session.Register("mongo", pder)
}
