package service

import (
	"blog_project/model"
	"blog_project/repo"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogService struct {
	BlogRepo           *repo.BlogRepository
	StakeholderService *StakeholderService
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

func (service *BlogService) HandleLike(ctx context.Context, blogId primitive.ObjectID, accountId string) (bool, error) {

	exists, err := service.StakeholderService.FindAccount(accountId)

	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("account does not exists")

	}

	liked, err := service.BlogRepo.HasLiked(ctx, blogId, accountId)
	if err != nil {
		return false, err
	}

	if liked {
		err = service.BlogRepo.RemoveLike(ctx, blogId, accountId)
		return false, err
	} else {
		err = service.BlogRepo.AddLike(ctx, blogId, accountId)
		return true, err
	}
}
