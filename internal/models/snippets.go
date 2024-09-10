package models

import (
	"context"
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Snippet struct {
	ID      string
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB      *sql.DB
	DBMongo *mongo.Client
}

func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {

	res, err := m.DBMongo.Database("snippets").Collection("snippets").InsertOne(context.TODO(), bson.D{
		{Key: "Title", Value: title},
		{Key: "Created", Value: time.Now().Local()},
		{Key: "Content", Value: content},
		{Key: "Expires", Value: expires},
	})
	if err != nil {
		return "0", err
	}
	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func (m *SnippetModel) Get(id string) (*Snippet, error) {
	var snippet Snippet
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := m.DBMongo.Database("snippets").Collection("snippets").FindOne(context.TODO(), bson.M{"_id": objectId})

	JSON := bson.M{}

	err = res.Decode(&JSON)

	if err != nil {
		return nil, err
	}

	expires := JSON["Expires"].(int32)
	snippet.Expires = time.Now().Local().AddDate(0, 0, int(expires))
	snippet.Title = JSON["Title"].(string)
	snippet.Content = JSON["Content"].(string)
	snippet.Created = JSON["Created"].(primitive.DateTime).Time()
	snippet.ID = JSON["_id"].(primitive.ObjectID).Hex()

	return &snippet, err
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	var snippets []*Snippet
	opts := options.Find().SetLimit(10)

	cursor, err := m.DBMongo.Database("snippets").Collection("snippets").Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, result := range results {
		var snippet Snippet
		expires := result["Expires"].(int32)
		snippet.Expires = time.Now().Local().AddDate(0, 0, int(expires))
		snippet.Title = result["Title"].(string)
		snippet.Content = result["Content"].(string)
		snippet.Created = result["Created"].(primitive.DateTime).Time()
		snippet.ID = result["_id"].(primitive.ObjectID).Hex()
		snippets = append(snippets, &snippet)
	}

	return snippets, nil
}
