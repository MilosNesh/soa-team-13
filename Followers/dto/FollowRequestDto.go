package dto

type FollowRequest struct {
	FollowerId string `json:"followerId"`
	FollowedId string `json:"followedId"`
}
