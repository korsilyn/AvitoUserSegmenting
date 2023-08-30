package v1

import (
	"avito-user-segmenting/core/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type slugRoutes struct {
	slugService service.Slug
}

func newSlugRoutes(g *echo.Group, slugService service.Slug) *slugRoutes {
	r := &slugRoutes{
		slugService: slugService,
	}

	g.POST("/create", r.create)
	g.POST("/remove", r.remove)

	return r
}

type slugInput struct {
	Name string `json:"name"`
}

func (r *slugRoutes)create(c echo.Context) error {
	var input slugInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	id, err := r.slugService.CreateSlug(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

func (r *slugRoutes)remove(c echo.Context) error {
	var input slugInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	err := r.slugService.RemoveSlugByName(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}
	
	type response struct {}

	return c.JSON(http.StatusOK, response{})
}
