package dto

type PostBlogDTO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageUrl string `json:"image_url,omitempty"`
}
