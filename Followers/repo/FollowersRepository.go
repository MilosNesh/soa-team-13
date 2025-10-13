package repo

import (
	"context"

	"followers.com/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowRepository struct {
	Driver neo4j.DriverWithContext
}

func (r *FollowRepository) Follow(follower, followed *model.Profile) error {
	session := r.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(context.TODO(),
			`MERGE (follower:Profile {username: $followerUsername})
			ON CREATE SET follower.id = $followerId
			MERGE (followed:Profile {username: $followedUsername})
			ON CREATE SET followed.id = $followedId
			MERGE (follower)-[:FOLLOWS]->(followed)
			`,
			map[string]interface{}{
				"followerId":       follower.Id,
				"followerUsername": follower.Username,
				"followedId":       followed.Id,
				"followedUsername": followed.Username,
			},
		)
		return nil, err
	})
	return err
}

func (r *FollowRepository) IsFollowing(followerId, followedId string) (bool, error) {
	session := r.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.TODO())

	result, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (follower:Profile {id: $followerId})-[:FOLLOWS]->(followed:Profile {id: $followedId})
			RETURN count(*) > 0 as isFollowing
		`
		result, err := tx.Run(context.TODO(), query, map[string]interface{}{
			"followerId": followerId,
			"followedId": followedId,
		})
		if err != nil {
			return false, err
		}

		if result.Next(context.TODO()) {
			return result.Record().Values[0].(bool), nil
		}
		return false, nil
	})

	if err != nil {
		return false, err
	}
	return result.(bool), nil
}
