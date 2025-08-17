package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DatabaseExists проверяет существование базы данных.
func DatabaseExists(client *mongo.Client, name string) (bool, error) {
	dbNamesList, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		return false, fmt.Errorf("data bases names request err:%w", err)
	}

	for _, dbName := range dbNamesList {
		if dbName == name {
			return true, nil
		}
	}

	return false, nil
}

// DatabaseCreate создает базу данных.
func DatabaseCreate(client *mongo.Client, name string) error {
	// в MongoDB база создается автоматически при первой записи
	db := client.Database(name)
	if err := db.CreateCollection(context.Background(), "temp_collection"); err != nil {
		return fmt.Errorf("db create err: %w", err)
	}

	// Удаляем временную коллекцию, так как база создана.
	if err := db.Collection("temp_collection").Drop(context.Background()); err != nil {
		return fmt.Errorf("temporary collection clean up err: %w", err)
	}

	return nil
}

func DatabaseDrop(client *mongo.Client, name string) error {
	if err := client.Database(name).Drop(context.Background()); err != nil {
		return fmt.Errorf("db drop err: %w", err)
	}

	return nil
}

// CollectionList возвращает список коллекций в базе данных.
func CollectionList(client *mongo.Client, dbName string) ([]string, error) {
	collections, err := client.Database(dbName).ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("list collections names request err: %w", err)
	}

	return collections, nil
}

// CollectionDrop удаляет указанные коллекции в базе данных.
func CollectionDrop(client *mongo.Client, dbName string, names ...string) error {
	db := client.Database(dbName)
	for _, name := range names {
		err := db.Collection(name).Drop(context.Background())
		if err != nil {
			return fmt.Errorf("collection %s drop err: %w", name, err)
		}
	}

	return nil
}

// DatabaseClear удаляет все коллекции из базы данных.
func DatabaseClear(client *mongo.Client, dbName string) error {
	collections, err := CollectionList(client, dbName)
	if err != nil {
		return err
	}

	err = CollectionDrop(client, dbName, collections...)
	if err != nil {
		return err
	}

	return nil
}
