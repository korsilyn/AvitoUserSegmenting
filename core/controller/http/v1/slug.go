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

// @Summary Create slug
// @Description Create slug
// @Tags slugs
// @Accept json
// @Produce json
// @Success 201 {object} v1.slugRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/slugs/create [post]
func (r *slugRoutes)create(c echo.Context) error {
	var input slugInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	id, err := r.slugService.CreateSlug(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

// @Summary Remove slug
// @Description Remove slug
// @Tags slugs
// @Accept json
// @Produce json
// @Success 201 {object} v1.slugRoutes.remove.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/slugs/remove [post]
func (r *slugRoutes)remove(c echo.Context) error {
	var input slugInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	err := r.slugService.RemoveSlugByName(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	
	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusOK, response{
		Message: "Success",
	})
}
