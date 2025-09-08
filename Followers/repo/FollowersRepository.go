package repo

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowRepository struct {
	Driver neo4j.DriverWithContext
}

func (r *FollowRepository) Follow(followerId, followedId string) error {
	session := r.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(context.TODO(),
			`MERGE (follower:Profile {id: $followerId})
             MERGE (followed:Profile {id: $followedId})
             MERGE (follower)-[:FOLLOWS]->(followed)`,
			map[string]interface{}{"followerId": followerId, "followedId": followedId},
		)
		return nil, err
	})
	return err
}
