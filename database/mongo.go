package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"random_image/config"
	"random_image/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB struct
type MongoDB struct {
	ic *mongo.Collection
}

// ImageObj struct
type ImageObj struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ImageID     int                `bson:"image_id" json:"id"`
	ImageDetail model.Image        `bson:"detail" json:"details"`
}

// ErrUserNotFound is returned if user with credentials is not found
var ErrUserNotFound = errors.New("user not found")

// ErrInvestmentNotFound is returned if user with credentials is not found
var ErrInvestmentNotFound = errors.New("Investment not found")

// NewMongoDB function
func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://random_images:%s@cluster0.4ykaj.mongodb.net/%s?retryWrites=true&w=majority", cfg.DatastoreDBPassword, cfg.DatastoreDBName)))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	ic := client.Database("random_images").Collection("images")
	return &MongoDB{ic: ic}, nil
}

// AddNewImage method
func (m *MongoDB) AddNewImage(id int, img *model.Image) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	newObj := &ImageObj{
		ImageID:     id,
		ImageDetail: *img,
	}
	_, err := m.ic.InsertOne(ctx, newObj)
	return err
}

// GetImageByID method
func (m *MongoDB) GetImageByID(id int) (*model.Image, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	f := bson.M{"image_id": id}
	result := m.ic.FindOne(ctx, f)
	var obj ImageObj
	err := result.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj.ImageDetail, nil
}

// GetAll method
func (m *MongoDB) GetAll() ([]model.Image, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	f := bson.M{}
	result, err := m.ic.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var obj ImageObj
	res := []model.Image{}
	for result.Next(context.TODO()) {
		err = result.Decode(&obj)
		if err != nil {
			return nil, err
		}
		res = append(res, obj.ImageDetail)
	}
	return res, nil
}
