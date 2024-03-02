package providers

import (
	"container/list"
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/eyko139/go-snippets/internal/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pder = &MongoSessionProvider{sessionLogger: log.New(os.Stdout, "Session\t", log.Ldate|log.Ltime)}

type MongoSessionStore struct {
	sid          string                 `bson:"sid"`
	timeAccessed time.Time              `bson:"timeAccessed"`
	value        map[string]interface{} `bson:"value"`
}

type MongoSessionProvider struct {
	lock          sync.Mutex
	collection    *mongo.Collection
	sessions      map[string]*list.Element
	sessionLogger *log.Logger
}

func (mss *MongoSessionStore) Set(key string, value interface{}) error {
	mss.value[key] = value
	pder.SessionUpdate(mss.sid, mss)
	return nil
}

func (mss *MongoSessionStore) Get(key string) interface{} {
	return mss.value[key]
}

func (mss *MongoSessionStore) Delete(key string) error {
	_, err := pder.collection.DeleteOne(context.TODO(), key)
	return err
}

func (mss *MongoSessionStore) SessionID() string {
	return ""
}

func (msp *MongoSessionProvider) SessionInit(sid string) (session.Session, error) {
	v := make(map[string]interface{})
	newSession := &MongoSessionStore{sid: sid, timeAccessed: time.Now().Local(), value: v}
	_, err := msp.collection.InsertOne(context.TODO(), bson.D{
        {Key: "sid", Value: newSession.sid}, 
        {Key: "timeAccessed", Value: newSession.timeAccessed}, 
        {Key: "value", Value: v},
    })
	if err != nil {
		panic(err)
	}
	return newSession, nil
}

type myTime time.Time

func (msp *MongoSessionProvider) SessionRead(sid string) (session.Session, error) {
	var result map[string]interface{}
	err := msp.collection.FindOne(context.TODO(), bson.D{{Key: "sid", Value: sid}}).Decode(&result)
	timeAccessed := result["timeAccessed"].(primitive.DateTime).Time()
	if err != nil {
		panic(err)
	}
	value, ok := result["value"].(map[string]interface{})
	if !ok {
		panic("could not cast value")
	}
	session := &MongoSessionStore{
		sid:          result["sid"].(string),
		timeAccessed: timeAccessed,
		value:        value,
	}
	return session, nil
}

func (msp *MongoSessionProvider) SessionDestroy(sid string) error {
	_, err := msp.collection.DeleteOne(context.TODO(), sid)
	return err
}
func (msp *MongoSessionProvider) SessionGC(maxLifeTime int64) {
	threshold := primitive.NewDateTimeFromTime(time.Now().Add(-time.Second * time.Duration(maxLifeTime)))
	msp.sessionLogger.Printf("Deleting Sessions older than %s\n", threshold.Time())
	filter := bson.M{"timeAccessed": bson.M{"$lt": threshold}}
	deleteResult, err := msp.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	msp.sessionLogger.Printf("Deleted %d documents\n", deleteResult.DeletedCount)
}

func (msp *MongoSessionProvider) SessionUpdate(sid string, update *MongoSessionStore) error {
	_, err := msp.collection.UpdateOne(
		context.TODO(),
		bson.D{{Key: "sid", Value: sid}},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "sid", Value: sid},
				{Key: "timeAccessed", Value: time.Now()},
				{Key: "value", Value: update.value},
			},
        }})
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
