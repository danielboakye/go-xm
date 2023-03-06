package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/models"
	"github.com/danielboakye/go-xm/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	r   repo.IRepository
	v   *helpers.Validation
	cfg config.Configurations
}

func NewHandler(r repo.IRepository, v *helpers.Validation, c config.Configurations) *Handler {
	return &Handler{r: r, v: v, cfg: c}
}

func (h *Handler) companyExists(ctx context.Context, companyName string) (exists bool, err error) {
	_, err = h.r.GetCompanyByName(ctx, companyName)
	if err == sql.ErrNoRows {
		err = nil
		return
	}

	exists = err == nil
	return
}

func (h *Handler) CreateCompany(c *gin.Context) {

	var request models.Company
	if err := c.ShouldBind(&request); err != nil {
		err = helpers.ErrInvalidParameters
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := h.v.ValidateForm(request); err != nil {
		err = helpers.ErrInvalidParameters
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	exists, err := h.companyExists(c.Request.Context(), request.Name)
	if err != nil {
		log.Println(err)
		err = helpers.ErrProcessingFailed
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	if exists {
		err = helpers.ErrDuplicateRecord
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	companyID, err := h.r.CreateCompany(
		c.Request.Context(),
		request.Name,
		request.Description,
		request.AmountOfEmployees,
		request.Registered,
		request.CompanyType,
	)

	if err != nil {
		err = helpers.ErrProcessingFailed
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	request.ID = companyID

	c.JSON(http.StatusOK, request)
}

func (h *Handler) UpdateCompany(c *gin.Context) {

	var request models.CompanyUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		err = helpers.ErrInvalidParameters
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	if err := h.v.ValidateForm(request); err != nil {
		log.Println(err)
		err = helpers.ErrInvalidParameters
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	companyIDParam := c.Param("company-id")

	data, err := h.r.GetCompanyByID(
		c.Request.Context(),
		companyIDParam,
	)

	if err != nil {
		log.Println(err)
		err = helpers.ErrProcessingFailed
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	log.Println("here1")
	if request.Name != nil {
		if *request.Name != data.Name {
			exists, err := h.companyExists(c.Request.Context(), *request.Name)
			if err != nil {
				err = helpers.ErrProcessingFailed
				c.AbortWithStatusJSON(
					helpers.GetHttpStatusByErr(err),
					gin.H{"error": err.Error()},
				)
				return
			}

			if exists {
				err = helpers.ErrDuplicateRecord
				c.AbortWithStatusJSON(
					helpers.GetHttpStatusByErr(err),
					gin.H{"error": err.Error()},
				)
				return
			}
		}
	}

	log.Println("here2")
	err = h.r.UpdateCompany(
		c.Request.Context(),
		companyIDParam,
		request.Name,
		request.Description,
		request.AmountOfEmployees,
		request.Registered,
		request.CompanyType,
	)

	log.Println("here3")
	if err != nil {
		log.Println(err)
		err = helpers.ErrProcessingFailed
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteCompany(c *gin.Context) {

	companyIDParam := c.Param("company-id")

	err := h.r.DeleteCompany(
		c.Request.Context(),
		companyIDParam,
	)

	if err != nil {
		err = helpers.ErrProcessingFailed
		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetCompany(c *gin.Context) {

	companyIDParam := c.Param("company-id")

	data, err := h.r.GetCompanyByID(
		c.Request.Context(),
		companyIDParam,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			err = helpers.ErrNoRecordFound
		} else {
			err = helpers.ErrProcessingFailed
		}

		c.AbortWithStatusJSON(
			helpers.GetHttpStatusByErr(err),
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, data)
}
