package metadata_sqlite

import (
	"fmt"
	"log"

	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

type FaceEmbeddingQueryBuilder struct {
}

func (b FaceEmbeddingQueryBuilder) GetFindSimilarQuery(ctx *metadata_db_entity.FacesContext, embedding *[]float32, max int) string {
	return "WITH knn_matches AS (SELECT id, assetId, distance FROM faces WHERE embedding MATCH vec_f32(?) AND k = 99)" +
		"SELECT id, assetId, distance FROM knn_matches WHERE distance <= 25.0 ORDER BY distance;"
}

type FaceEmbeddingSelectQueryBuilder struct {
	defaultSelectQueryBuilder go_db_entities.SelectQueryBuilder
}

func (b FaceEmbeddingSelectQueryBuilder) GetSelectQuery(c *go_db_entities.Context, e go_db_entities.Entity, filter go_db_entities.Filter) string {

	if _, isFaceEmbedding := e.(*metadata_db_entity.FaceEmbedding); isFaceEmbedding {
		sql := "SELECT Id, AssetId, X1, Y1, X2, Y2, vec_to_json(Embedding) FROM Faces"

		where := c.QueryBuilder.FilterQueryBuilder.GetFilterStatement(c, filter)
		if len(where) > 0 {
			sql += " WHERE " + where
		}

		if filter.OrderBy != nil {
			order := ""
			for _, field := range filter.OrderBy {
				if len(order) > 0 {
					order += ","
				}
				order += field
			}
			if len(order) > 0 {
				sql += " ORDER BY " + order
			}
		}

		if filter.Offset > 0 || filter.Count > 0 {
			sql += fmt.Sprintf(" LIMIT %d,%d", filter.Offset, filter.Count)
		}

		if filter.Next != nil && filter.Next.Offset > 0 && filter.Next.Count > 0 {
			log.Printf("Warn: Ignoring offest/count of nested filter")
		}

		sql += ";"

		return sql
	}

	return b.defaultSelectQueryBuilder.GetSelectQuery(c, e, filter)
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
