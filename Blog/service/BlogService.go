package service

import (
	"blog_project/model"
	"blog_project/repo"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogService struct {
	BlogRepo *repo.BlogRepository
}

func (service *BlogService) FindBlogById(ctx context.Context, id string) (*model.Blog, error) {
	blog, err := service.BlogRepo.FindById(ctx, id)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}

		return nil, err
	}
	return &blog, nil
}

func (service *BlogService) Create(ctx context.Context, blog *model.Blog) error {
	err := service.BlogRepo.Create(ctx, blog)
	if err != nil {
		return err
	}
	return nil
}

func (service *BlogService) HandleLike(ctx context.Context, blogId primitive.ObjectID, username string) (bool, error) {

	liked, err := service.BlogRepo.HasLiked(ctx, blogId, username)
	if err != nil {
		return false, err
	}

	if liked {
		err = service.BlogRepo.RemoveLike(ctx, blogId, username)
		return false, err
	} else {
		err = service.BlogRepo.AddLike(ctx, blogId, username)
		return true, err
	}
}
