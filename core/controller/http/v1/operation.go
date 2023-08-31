package v1

import (
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
	Percent int `json:"percent"`
	TTL int `json:"ttl"`
}

// @Summary Create operation
// @Description Add slugs to some percent of users
// @Tags operations
// @Accept json
// @Produce json
// @Param input body addSlugsInput true "input"
// @Success 200 {object} v1.operationRoutes.addSlugs.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/operations/add [post]
func (r *operationRoutes) addSlugs(c echo.Context) error {
	var input addSlugsInput

	if err := c.Bind(&input); err != nil || input.Percent < 0 || input.Percent > 100 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return err
	}

	err := r.operationService.CreateOperations(c.Request().Context(), service.OperationCreateInput{
		Slugs: input.Slugs,
		Percent: input.Percent,
		TTL: input.TTL,
	})
	if err != nil {
		log.Debugf("contoller.addSlugsInput - CreateOperations: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return err
	}

	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusOK, response{
		Message: "Successful",
	})
}

type removeSlugsInput struct {
	Slugs []string `json:"slugs"`
	UserId int `json:"user_id"`
}

// @Summary Remove operation
// @Description Remove slug from user
// @Tags operations
// @Accept json
// @Produce json
// @Param input body removeSlugsInput true "input"
// @Success 200 {object} v1.operationRoutes.removeSlugs.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/operations/remove [post]
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

	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusOK, response{
		Message: "Successful",
	})
}

type getSlugsInput struct {
	UserId int `json:"user_id"`
}

// @Summary Get operations
// @Description Get all slugs of user
// @Tags operations
// @Accept json
// @Produce json
// @Param input body getSlugsInput true "input"
// @Success 200 {object} v1.operationRoutes.getSlugs.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/operations/ [get]
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
		Slugs []string `json:"slugs"`
	}

	return c.JSON(http.StatusOK, response{
		Slugs: output,
	})
}
