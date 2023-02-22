package db_client

import (
	"context"
	"errors"
	"task1/items_manager/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) InsertItem(item models.Item) (string, error) {
	item.Created = time.Now()
	item.Modified = time.Now()
	res, err := c.itemsCollection.InsertOne(context.Background(), item)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) GetAllUserItems(username string) ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cursor, err := c.itemsCollection.Find(ctx, bson.M{"owner": username})
	if err != nil {
		return []models.Item{}, err
	}
	var items []models.Item

	if err = cursor.All(ctx, &items); err != nil {
		panic(err)
	}

	for i := range items {
		items[i].Id = items[i].PrimitiveId.Hex()
	}
	return items, nil
}

func (c *Client) GetItem(id string) (models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	value := c.itemsCollection.FindOne(ctx, bson.M{"_id": objID})
	var item models.Item
	if errors.Is(value.Err(), mongo.ErrNoDocuments) {
		return item, models.ErrItemDoesNotExist
	}
	value.Decode(&item)
	item.Id = item.PrimitiveId.Hex()
	return item, nil
}

func (c *Client) DeleteItem(item models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	obj, err := primitive.ObjectIDFromHex(item.Id)

	if err != nil {
		return models.ErrDatabaseOperation
	}
	item.PrimitiveId = obj
	del, err := c.itemsCollection.DeleteOne(ctx, bson.M{"_id": item.PrimitiveId})

	if err != nil {
		panic(err)
	}

	if del.DeletedCount == 0 {
		return models.ErrItemDoesNotExist
	}
	return nil
}

func (c *Client) UpdateItem(item *models.UpdateItem) error {
	oldItem, err := c.GetItem(item.Id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	updatedItem := item.Update(oldItem)
	id, _ := primitive.ObjectIDFromHex(item.Id)

	_, err = c.itemsCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		updatedItem,
	)
	// item.PrimitiveId = obj
	// // filter :
	// upd, err := c.itemsCollection.UpdateByID(ctx, obj, bson.D{{"$set":item}})
	if err != nil {
		return err
	}
	return nil
}
