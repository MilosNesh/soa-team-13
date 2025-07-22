package service

import (
	"blog_project/model"
	"blog_project/repo"
	"context"

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
