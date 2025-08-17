package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/wal1251/pkg/core/errs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Storage = (*StorageImpl)(nil)

// StorageImpl — реализация интерфейса Storage.
// Отвечает за взаимодействие с базой данных MongoDB.
type StorageImpl struct {
	DB *mongo.Database // Экземпляр базы данных MongoDB.
}

// Storage определяет интерфейс для выполнения основных операций с базой данных.
// Этот интерфейс позволяет гибко работать с разными реализациями хранилища.
type Storage interface {
	// InsertOne вставляет один документ в указанную коллекцию.
	// Возвращает ID вставленного документа или ошибку в случае сбоя.
	InsertOne(ctx context.Context, collection string, document Document) (ID, error)

	// InsertMany вставляет несколько документов в указанную коллекцию.
	// Возвращает список ID вставленных документов или ошибку в случае сбоя.
	InsertMany(ctx context.Context, collection string, documents []Document) ([]ID, error)

	// DeleteOne удаляет один документ из указанной коллекции по ID.
	// Возвращает количество удалённых документов (0 или 1) или ошибку в случае сбоя.
	DeleteOne(ctx context.Context, collection string, id ID) (int64, error)

	// DeleteMany удаляет несколько документов из указанной коллекции по списку ID.
	// Возвращает количество удалённых документов или ошибку в случае сбоя.
	DeleteMany(ctx context.Context, collection string, ids []ID) (int64, error)

	// FindOne находит один документ в указанной коллекции по фильтру и проекции.
	// Возвращает найденный документ или ошибку, если документ не найден.
	FindOne(ctx context.Context, collection string, filter *Filter, projection *Projection) (*Document, error)

	// FindMany находит несколько документов в указанной коллекции по фильтру и проекции.
	// Возвращает список найденных документов или ошибку в случае сбоя.
	FindMany(ctx context.Context, collection string, filter *Filter, projection *Projection) ([]*Document, error)

	// UpdateOne обновляет один документ в указанной коллекции, соответствующий фильтру.
	// Возвращает количество обновлённых документов (0 или 1) или ошибку в случае сбоя.
	UpdateOne(ctx context.Context, collection string, filter *Filter, update *Update) (int64, error)

	// UpdateMany обновляет несколько документов в указанной коллекции, соответствующих фильтру.
	// Возвращает количество обновлённых документов или ошибку в случае сбоя.
	UpdateMany(ctx context.Context, collection string, filter *Filter, update *Update) (int64, error)

	// UpdateByID обновляет один документ в указанной коллекции по его ID.
	// Возвращает количество обновлённых документов (0 или 1) или ошибку в случае сбоя.
	UpdateByID(ctx context.Context, collection string, id ID, update *Update) (int64, error)

	// ReplaceOne заменяет один документ в указанной коллекции, соответствующий фильтру.
	// Возвращает количество обновлённых документов (0 или 1) или ошибку в случае сбоя.
	ReplaceOne(ctx context.Context, collection string, filter *Filter, document Document) (int64, error)

	// Disconnect отключает соединение с базой данных.
	Disconnect(ctx context.Context) error

	// HealthCheck проверяет состояние базы данных.
	HealthCheck(ctx context.Context) error
}

// NewMongoStorage создаёт новый экземпляр StorageImpl для работы с базой данных MongoDB.
// Принимает конфигурацию для подключения к MongoDB.
// Возвращает экземпляр StorageImpl или ошибку, если подключение не удалось.
func NewMongoStorage(ctx context.Context, config *Config) (*StorageImpl, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetURI()))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к MongoDB: %w", err)
	}

	db := client.Database(config.Database)

	return &StorageImpl{
		DB: db,
	}, nil
}

func (s *StorageImpl) InsertOne(ctx context.Context, collection string, document Document) (ID, error) {
	result, err := s.DB.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		return ID{}, fmt.Errorf("failed to insert document: %w", err)
	}

	return ID{value: result.InsertedID.(primitive.ObjectID)}, nil //nolint:forcetypeassert
}

func (s *StorageImpl) InsertMany(ctx context.Context, collection string, documents []Document) ([]ID, error) {
	anyList := make([]any, 0, len(documents))
	for _, e := range documents {
		anyList = append(anyList, e)
	}

	result, err := s.DB.Collection(collection).InsertMany(ctx, anyList)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	idList := make([]ID, 0, len(result.InsertedIDs))
	for _, r := range result.InsertedIDs {
		idList = append(idList, ID{value: r.(primitive.ObjectID)}) //nolint:forcetypeassert
	}

	return idList, nil
}

func (s *StorageImpl) DeleteOne(ctx context.Context, collection string, id ID) (int64, error) {
	result, err := s.DB.Collection(collection).DeleteOne(ctx, bson.M{"_id": id.value})
	if err != nil {
		return 0, fmt.Errorf("failed to delete document: %w", err)
	}

	return result.DeletedCount, nil
}

func (s *StorageImpl) DeleteMany(ctx context.Context, collection string, ids []ID) (int64, error) {
	result, err := s.DB.Collection(collection).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}

func (s *StorageImpl) FindOne(ctx context.Context, collection string, filter *Filter, projection *Projection) (*Document, error) {
	coll := s.DB.Collection(collection)

	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection.Value())
	}

	var document Document
	err := coll.FindOne(ctx, filter.Value(), opts).Decode(&document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.ErrNotFound
		}

		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	return &document, nil
}

func (s *StorageImpl) FindMany(ctx context.Context, collection string, filter *Filter, projection *Projection) ([]*Document, error) {
	coll := s.DB.Collection(collection)

	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection.Value())
	}

	cur, err := coll.Find(ctx, filter.Value(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	var documents []*Document
	if err = cur.All(ctx, &documents); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return documents, nil
}

func (s *StorageImpl) UpdateOne(ctx context.Context, collection string, filter *Filter, update *Update) (int64, error) {
	coll := s.DB.Collection(collection)

	result, err := coll.UpdateOne(ctx, filter.Value(), update.Value())
	if err != nil {
		return 0, fmt.Errorf("failed to update document: %w", err)
	}

	return result.ModifiedCount, nil
}

func (s *StorageImpl) UpdateMany(ctx context.Context, collection string, filter *Filter, update *Update) (int64, error) {
	coll := s.DB.Collection(collection)

	result, err := coll.UpdateMany(ctx, filter.Value(), update.Value())
	if err != nil {
		return 0, fmt.Errorf("failed to update documents: %w", err)
	}

	return result.ModifiedCount, nil
}

func (s *StorageImpl) UpdateByID(ctx context.Context, collection string, id ID, update *Update) (int64, error) {
	coll := s.DB.Collection(collection)

	result, err := coll.UpdateByID(ctx, id, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update document: %w", err)
	}

	return result.ModifiedCount, nil
}

func (s *StorageImpl) ReplaceOne(ctx context.Context, collection string, filter *Filter, document Document) (int64, error) {
	coll := s.DB.Collection(collection)

	result, err := coll.ReplaceOne(ctx, filter, document)
	if err != nil {
		return 0, fmt.Errorf("failed to replace document: %w", err)
	}

	return result.ModifiedCount, nil
}

func (s *StorageImpl) Disconnect(ctx context.Context) error {
	if err := s.DB.Client().Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect mongo client: %w", err)
	}

	return nil
}

func (s *StorageImpl) HealthCheck(ctx context.Context) error {
	if err := s.DB.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping mongo client: %w", err)
	}

	return nil
}
