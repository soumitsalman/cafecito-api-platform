package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soumitsalman/cafecito-api-platform/apis/espresso/db"
	utils "github.com/soumitsalman/cafecito-api-platform/apis/shared"
	datautils "github.com/soumitsalman/data-utils"
)

// These routes are for internal use ONLY.
// They are not part of the public API.
// Exclude them from swagger docs and `espresso.oas.json` definitions.

// previewConfidenceParams is the GET /preview/confidence request.
type previewConfidenceParams struct {
	IDs []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"required,min=1,max=128"`
}

// previewConfidenceDocument is one requested id and the confidence of its newest derived signal.
// Confidence is JSON null when that signal has no confidence, or when no signal is derived from the id.
type previewConfidenceDocument struct {
	ID         uuid.UUID `json:"id"`
	Confidence *string   `json:"confidence"`
}

func (params *previewConfidenceParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return utils.NewAPIError(utils.API_ERROR_INVALID_REQUEST, err.Error())
	}
	return nil
}

func (r *Configuration) previewGetConfidence(c *gin.Context) {
	var params previewConfidenceParams
	if err := params.shouldBind(c); err != nil {
		writeError(c, err)
		return
	}
	rows, err := r.DB.QueryDerivedConfidence(c.Request.Context(), params.IDs)
	if err != nil {
		writeError(c, utils.NewAPIError(utils.API_ERROR_DB_ERROR, API_ERROR_MSG_OUR_BAD))
		return
	}
	docs := datautils.Transform(rows, func(row *db.IDConfidence) previewConfidenceDocument {
		return previewConfidenceDocument{ID: row.ID, Confidence: row.Confidence}
	})
	writePage(c, docs, len(docs), nil, "json")
}
