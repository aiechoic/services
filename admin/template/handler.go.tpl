package {{.PackageName}}

import (
	"github.com/aiechoic/services/admin/internal/rsp"
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	db *DB
}

func NewHandlers(db *DB) *Handlers {
	return &Handlers{
		db: db,
	}
}

func (h *Handlers) Get() gins.Handler {
	type getRequest struct {
		By    string `form:"by" binding:"required"`
		Value string `form:"value" binding:"required"`
	}
	return gins.Handler{
		Request: gins.Request{
			Json: &getRequest{},
		},
		Response: gins.Response{
			Json: &{{.ModelName}}{},
		},
		Handler: func(c *gin.Context) {
			var req getRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			model, err := h.db.GetByColumn(req.By, req.Value)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", model))
		},
	}
}

func (h *Handlers) Create() gins.Handler {
	return gins.Handler{
		Request: gins.Request{
			Json: &{{.ModelName}}{},
		},
		Response: gins.Response{
			Json: &{{.ModelName}}{},
		},
		Handler: func(c *gin.Context) {
			var model {{.ModelName}}
			if err := c.BindJSON(&model); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			err := h.db.Create(&model)
			if err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", model))
		},
	}
}

func (h *Handlers) FullUpdate() gins.Handler {
	return gins.Handler{
		Request: gins.Request{
			Json: &{{.ModelName}}{},
		},
		Response: gins.Response{
			Json: &{{.ModelName}}{},
		},
		Handler: func(c *gin.Context) {
			var req {{.ModelName}}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			model, err := h.db.FullUpdateByColumn(&req, "id", req.ID)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", model))
		},
	}
}

func (h *Handlers) PartialUpdate() gins.Handler {
	return gins.Handler{
		Request: gins.Request{
			Json: &{{.ModelName}}{},
		},
		Response: gins.Response{
			Json: &{{.ModelName}}{},
		},
		Handler: func(c *gin.Context) {
			var req {{.ModelName}}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			model, err := h.db.PartialUpdateByColumn(&req, "id", req.ID)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", model))
		},
	}
}

func (h *Handlers) Delete() gins.Handler {
	type deleteRequest struct {
		By    string `json:"by" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	type deleteResponse struct {
		RowsAffected int64 `json:"rows_affected"`
	}
	return gins.Handler{
		Request: gins.Request{
			Json: &deleteRequest{},
		},
		Response: gins.Response{
			Json: &deleteResponse{},
		},
		Handler: func(c *gin.Context) {
			var req deleteRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			rowsAffected := h.db.DeleteByColumn(req.By, req.Value)
			c.JSON(200, rsp.Success("ok", &deleteResponse{RowsAffected: rowsAffected}))
		},
	}
}

func (h *Handlers) Count() gins.Handler {
	type countResponse struct {
		Count int64 `json:"count"`
	}
	return gins.Handler{
		Response: gins.Response{
			Json: &countResponse{},
		},
		Handler: func(c *gin.Context) {
			count, err := h.db.Count()
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", &countResponse{Count: count}))
		},
	}
}

func (h *Handlers) List() gins.Handler {
	type listRequest struct {
		By     string `json:"by"`
		Desc   bool   `json:"desc"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}
	return gins.Handler{
		Request: gins.Request{
			Json: &listRequest{},
		},
		Response: gins.Response{
			Json: []{{.ModelName}}{},
		},
		Handler: func(c *gin.Context) {
			var req listRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(200, rsp.BadRequestError(err))
				return
			}
			models, err := h.db.FindByColumn(req.By, req.Desc, req.Limit, req.Offset)
			if err != nil {
				c.JSON(200, rsp.InternalServerError(err))
				return
			}
			c.JSON(200, rsp.Success("ok", models))
		},
	}
}