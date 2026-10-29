package metadata_sqlite

import (
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

type FaceEmbeddingQueryBuilder struct {
}

func (b FaceEmbeddingQueryBuilder) GetFindSimilarQuery(ctx *metadata_db_entity.FacesContext, embedding *[]float32, max int) string {
	return "WITH knn_matches AS (SELECT id, assetId, distance FROM faces WHERE embedding MATCH vec_f32(?) AND k = 99)" +
		"SELECT id, assetId, distance FROM knn_matches WHERE distance <= 25.0 ORDER BY distance;"
}

type FaceEmbeddingCreateQueryBuilder struct {
	defaultCreateQueryBuilder go_db_entities.CreateQueryBuilder
}

func (b FaceEmbeddingCreateQueryBuilder) GetCreateQueries(c *go_db_entities.Context, e go_db_entities.Entity) []string {

	if _, isFaceEmbedding := e.(*metadata_db_entity.FaceEmbedding); isFaceEmbedding {
		return []string{
			"CREATE VIRTUAL TABLE IF NOT EXISTS faces USING vec0(id integer PRIMARY KEY, assetId integer, " +
				"x1 integer, y1 integer, x2 integer, y2 integer, " +
				"embedding FLOAT[512]);",
		}
	}

	return b.defaultCreateQueryBuilder.GetCreateQueries(c, e)
}
