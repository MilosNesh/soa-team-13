package repo

import (
	"blog_project/model"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	BlogsCollection    *mongo.Collection
	CommentsCollection *mongo.Collection
}

func (repo *BlogRepository) FindById(ctx context.Context, id string) (model.Blog, error) {
	blog := model.Blog{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return blog, mongo.ErrNoDocuments
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": objectId}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "blog_id",
			"as":           "comments",
		}},
		{"$limit": 1},
	}

	cursor, err := repo.BlogsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error during aggregation for blog ID %s: %v", id, err)
		return blog, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Printf("Blog with ID %s not found.", id)
		return blog, mongo.ErrNoDocuments
	}

	if err := cursor.Decode(&blog); err != nil {
		log.Printf("Error decoding blog result for ID %s: %v", id, err)
		return blog, err
	}

	if blog.Comments == nil {
		blog.Comments = []model.Comment{}
	}

	return blog, nil
}

func (repo *BlogRepository) Create(ctx context.Context, blog *model.Blog) error {
	insertResult, err := repo.BlogsCollection.InsertOne(ctx, blog)
	if err != nil {
		log.Printf("Greška prilikom kreiranja blog posta: %v", err)
		return err
	}
	if blog.Comments == nil {
		blog.Comments = []model.Comment{}
	}

	log.Printf("Uspešno kreiran blog post sa ID: %v", insertResult.InsertedID)

	if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		blog.ID = oid
	}

	return nil
}

func (repo *BlogRepository) FindAll(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog

	cursor, err := repo.BlogsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repo *BlogRepository) AddLike(ctx context.Context, blogId primitive.ObjectID, accountId string) error {
	filter := bson.M{"_id": blogId}
	update := bson.M{"$addToSet": bson.M{"likes": accountId}}

	result, err := repo.BlogsCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("blog not found add")
	}
	return nil
}

func (repo *BlogRepository) RemoveLike(ctx context.Context, blogId primitive.ObjectID, accountId string) error {
	filter := bson.M{"_id": blogId}
	update := bson.M{"$pull": bson.M{"likes": accountId}}

	result, err := repo.BlogsCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("blog not found")
	}

	return nil
}

func (repo *BlogRepository) HasLiked(ctx context.Context, blogId primitive.ObjectID, accountId string) (bool, error) {
	filter := bson.M{"_id": blogId, "likes": accountId}

	count, err := repo.BlogsCollection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (repo *BlogRepository) AddComment(ctx context.Context, comment *model.Comment) error {
	commentsCollection := repo.CommentsCollection

	_, err := commentsCollection.InsertOne(ctx, comment)
	if err != nil {
		log.Printf("Greška prilikom kreiranja komentara: %v", err)
		return err
	}

	return nil
}
