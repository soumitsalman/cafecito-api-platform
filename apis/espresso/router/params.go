package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bindableParams interface {
	shouldBind(c *gin.Context) error
}

// base response parameters for all endpoints
type responseTypeParams struct {
	ResponseType string `form:"response_type,default=json" binding:"omitempty,oneof=json yaml toon"`
}

// Path UUID params (Gin v1.12+ TextUnmarshaler parser).
type pathIDParams struct {
	ID uuid.UUID `uri:"id,parser=encoding.TextUnmarshaler" binding:"required"`
}

// base parameters for GET /events/{id}, /signals/{id}, /sources/{id}.
type itemParams struct {
	pathIDParams
	responseTypeParams
}

// base cursor-pagination params for all list queries.
type paginationParams struct {
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Cursor string `form:"cursor"`
	responseTypeParams
}

type vectorSearchParams struct {
	Q   string  `form:"q" binding:"max=1024"`
	Acc float64 `form:"acc,default=0.5" binding:"min=0,max=1"`
}

// base search params for all action, event and signal queries.
type sipQueryParams struct {
	From            time.Time `form:"from" time_format:"2006-01-02"`
	To              time.Time `form:"to" time_format:"2006-01-02"`
	Tags            []string  `form:"tags" collection_format:"csv" binding:"max=128"`
	Entities        []string  `form:"entities" collection_format:"csv" binding:"max=128,dive,oneof=company people product"`
	Categories      []string  `form:"categories" collection_format:"csv" binding:"max=128"`
	Companies       []string  `form:"companies" collection_format:"csv" binding:"max=128"`
	People          []string  `form:"people" collection_format:"csv" binding:"max=128"`
	Products        []string  `form:"products" collection_format:"csv" binding:"max=128"`
	Regions         []string  `form:"regions" collection_format:"csv" binding:"max=128"`
	ImpactedDomains []string  `form:"impacted_domains" collection_format:"csv" binding:"max=128"`
	ImpactLevels    []string  `form:"impact_levels" collection_format:"csv" binding:"max=128"`
}

// eventSearchParams holds filter and search parameters for GET /events.
type EventSearchParams struct {
	IDs        []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	EventTypes []string    `form:"event_types" collection_format:"csv" binding:"max=128"`
	SourceIDs  []uuid.UUID `form:"source_ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	sipQueryParams
	vectorSearchParams
	paginationParams
}

// signalSearchParams holds filter and search parameters for GET /signals.
type SignalSearchParams struct {
	IDs []uuid.UUID `form:"ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	sipQueryParams
	vectorSearchParams
	paginationParams
}

// eventEvidenceParams holds parameters for GET /events/{id}/evidence.
type EventEvidenceParams struct {
	SourceIDs []uuid.UUID `form:"source_ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	pathIDParams
	sipQueryParams
	paginationParams
}

// eventSignalsParams holds parameters for GET /events/{id}/signals.
type EventSignalsParams struct {
	pathIDParams
	sipQueryParams
	paginationParams
}

// signalEventsParams holds parameters for GET /signals/{id}/events.
type SignalEventsParams struct {
	SourceIDs  []uuid.UUID `form:"source_ids,parser=encoding.TextUnmarshaler" collection_format:"csv" binding:"max=128"`
	EventTypes []string    `form:"event_types" collection_format:"csv" binding:"max=128"`
	pathIDParams
	sipQueryParams
	paginationParams
}

// sourcesParams holds parameters for GET /sources.
type SourcesParams struct {
	Q       string   `form:"q" binding:"max=1024"`
	Domains []string `form:"domains" collection_format:"csv" binding:"max=128"`
	paginationParams
}

// tagsParams holds parameters for GET /tags.
type TagsParams struct {
	Q        string   `form:"q" binding:"max=1024"`
	Resource []string `form:"resource" collection_format:"csv" binding:"max=128,dive,oneof=event signal"`
	paginationParams
}

type EntitiesParams struct {
	Q     string   `form:"q" binding:"max=1024"`
	Types []string `form:"types" collection_format:"csv" binding:"max=3,dive,oneof=company people product"`
	paginationParams
}

type DiscoveryParams struct {
	Q string `form:"q" binding:"max=1024"`
	paginationParams
}

func (params *itemParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindUri(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *EventSearchParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *SignalSearchParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *EventEvidenceParams) shouldBind(c *gin.Context) error {
	path_params := pathIDParams{}
	if err := c.ShouldBindUri(&path_params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	params.pathIDParams = path_params
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *EventSignalsParams) shouldBind(c *gin.Context) error {
	path_params := pathIDParams{}
	if err := c.ShouldBindUri(&path_params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	params.pathIDParams = path_params
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *SignalEventsParams) shouldBind(c *gin.Context) error {
	path_params := pathIDParams{}
	if err := c.ShouldBindUri(&path_params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	params.pathIDParams = path_params
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *SourcesParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *TagsParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *EntitiesParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}

func (params *DiscoveryParams) shouldBind(c *gin.Context) error {
	if err := c.ShouldBindQuery(params); err != nil {
		return APIError{Code: API_ERROR_INVALID_REQUEST, Message: err.Error()}
	}
	return nil
}
