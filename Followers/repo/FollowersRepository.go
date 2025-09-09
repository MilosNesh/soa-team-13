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
