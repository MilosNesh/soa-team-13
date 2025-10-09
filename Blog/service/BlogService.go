package service

import (
	"blog_project/model"
	"blog_project/repo"
	"context"
	"errors"
	"strings"
	"time"

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
	processImageUrl(&blog)

	return &blog, nil
}

func (service *BlogService) Create(ctx context.Context, blog *model.Blog) error {
	err := service.BlogRepo.Create(ctx, blog)
	if err != nil {
		return err
	}
	return nil
}

func (s *BlogService) FindAllBlogs(ctx context.Context) ([]model.Blog, error) {
	blogs, err := s.BlogRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range blogs {
		processImageUrl(&blogs[i])
	}

	return blogs, nil
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

func (service *BlogService) AddComment(ctx context.Context, comment *model.Comment) error {
	_, err := service.BlogRepo.FindById(ctx, comment.BlogId.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("blog not found")
		}
		return err
	}

	authorName, err := service.StakeholderService.GetAuthorUsername(comment.AuthorId)
	if err != nil {
		return err
	}
	if authorName == "" {
		return errors.New("author not found")
	}

	comment.AuthorName = authorName

	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}

	err = service.BlogRepo.AddComment(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}

func processImageUrl(blog *model.Blog) {
	if blog.ImageUrl != "" && !strings.HasPrefix(blog.ImageUrl, "/uploads/") {
		blog.ImageUrl = "/uploads/" + blog.ImageUrl
	}
}
