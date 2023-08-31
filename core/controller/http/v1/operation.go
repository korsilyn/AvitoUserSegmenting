package v1

import (
	"avito-user-segmenting/core/entity"
	"avito-user-segmenting/core/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type operationRoutes struct {
	operationService service.Operation
}

func newOperationRoutes(g *echo.Group, operationService service.Operation) *operationRoutes {
	r := &operationRoutes{
		operationService: operationService,
	}

	g.POST("/add", r.addSlugs)
	g.POST("/remove", r.removeSlugs)
	g.GET("/", r.getSlugs)

	return r
}

type addSlugsInput struct {
	Slugs []string `json:"slugs"`
	UserId int `json:"user_id"`
	TTL int `json:"ttl"`
}

func (r *operationRoutes) addSlugs(c echo.Context) error {
	var input addSlugsInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	id, err := r.operationService.CreateOperations(c.Request().Context(), service.OperationCreateInput{
		Slugs: input.Slugs,
		UserId: input.UserId,
		TTL: input.TTL,
	})
	if err != nil {
		log.Debugf("contoller.addSlugsInput - CreateOperations: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusOK, response{
		Id: id,
	})
}

type removeSlugsInput struct {
	Slugs []string `json:"slugs"`
	UserId int `json:"user_id"`
}

func (r *operationRoutes) removeSlugs(c echo.Context) error {
	var input removeSlugsInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	err := r.operationService.RemoveOperationsByUserId(c.Request().Context(), service.OperationRemoveInput{
		Slugs: input.Slugs,
		UserId: input.UserId,
	})
	if err != nil {
		log.Debugf("contoller.removeSlugsInput - RemoveOperationsByUserId: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return err
	}

	type response struct {}

	return c.JSON(http.StatusOK, response{})
}

type getSlugsInput struct {
	UserId int `json:"user_id"`
}

func (r *operationRoutes) getSlugs(c echo.Context) error {
	var input getSlugsInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	output, err := r.operationService.GetAllSlugsByUserId(c.Request().Context(), input.UserId)
	if err != nil {
		log.Debugf("contoller.getSlugs - GetAllSlugsByUserId: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return err
	}

	type response struct {
		slugs []string `json:"slugs"`
	}

	return c.JSON(http.StatusOK, response{
		slugs: output,
	})
}
