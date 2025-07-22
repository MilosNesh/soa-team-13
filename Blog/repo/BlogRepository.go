package repo

import (
	"blog_project/model"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	Collection *mongo.Collection
}

func (repo *BlogRepository) FindById(ctx context.Context, id string) (model.Blog, error) {
	blog := model.Blog{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return blog, mongo.ErrNoDocuments
	}

	filter := bson.M{"_id": objectId}

	err = repo.Collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Blog with ID %s not found.", id)
			return blog, mongo.ErrNoDocuments
		}
		log.Printf("Error finding blog by ID %s: %v", id, err)
		return blog, err
	}

	return blog, nil
}

func (repo *BlogRepository) Create(ctx context.Context, blog *model.Blog) error {
	insertResult, err := repo.Collection.InsertOne(ctx, blog)
	if err != nil {
		log.Printf("Greška prilikom kreiranja blog posta: %v", err)
		return err
	}

	log.Printf("Uspešno kreiran blog post sa ID: %v", insertResult.InsertedID)

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		blog.ID = oid
	}

	return nil
}
