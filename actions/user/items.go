package user

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"net/http"
)

func ListUserItems(c *gin.Context) {
	pg := paginate.New()
	db := actions.GetDBConnectionByContext(c)
	actionUser := actions.GetActionUserByContext(c)
	repository := repositories.NewRepository(db)

	c.JSON(http.StatusOK, pg.With(repository.ListUserItems(actionUser)).Request(c.Request).Response(&[]models.Item{}))
}
