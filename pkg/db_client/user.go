package db_client

import (
	"context"
	"errors"
	"task1/items_manager/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (c *Client) CreateUser(user models.User) error {
	_, err := c.CheckUserNameExists(user.Username)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	ctx := context.Background()

	_, err = c.usersCollection.InsertOne(ctx, bson.D{
		{Key: "firstName", Value: user.FirstName},
		{Key: "lastName", Value: user.LastName},
		{Key: "username", Value: user.Username},
		{Key: "password", Value: string(hashedPassword)},
	})

	if err != nil {
		return err
	}

	return nil
}

// Results in an error if time exceeds or there is problem connecting
// to the database
func (c *Client) CheckUserNameExists(username string) (bool, error) {

	res := c.getUserDetails(username)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return false, nil
	}

	if res.Err() != nil {
		return false, models.ErrDatabaseOperation
	}

	return true, nil
}

func (c *Client) getUserDetails(username string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return c.usersCollection.FindOne(ctx, bson.M{"username": username})
}

func (c *Client) ValiatePassword(user models.User) error {
	userRes := c.getUserDetails(user.Username)

	if errors.Is(userRes.Err(), mongo.ErrNoDocuments) {
		return models.ErrUserDoesNotExist
	}
	var userInfo models.User
	userRes.Decode(&userInfo)
	err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(user.Password))

	if err != nil {
		return models.ErrInvalidPassword
	}

	return nil
}

func (c *Client) DeleteUser(user models.User) error {

	err := c.ValiatePassword(user)
	if err != nil {
		return err
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.usersCollection.DeleteOne(ctx, bson.M{"username": user.Username})

	if err != nil {
		panic(err)
	}

	if res.DeletedCount == 0 {
		return models.ErrUserDoesNotExist
	}

	return nil
}
