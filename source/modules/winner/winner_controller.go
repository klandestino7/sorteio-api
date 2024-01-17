package winner

import (
	"net/http"
	"sorteio-api/source/dto"

	"github.com/gin-gonic/gin"
)

// CONTROLLER
type IWinnerController interface {
	RequestAllWinners(c *gin.Context)
}

type WinnerController struct {
	WinnerService IWinnerService
}

func InitWinnerController(WinnerService IWinnerService) IWinnerController {
	return &WinnerController{
		WinnerService: WinnerService,
	}
}

func (s *WinnerController) RequestAllWinners(c *gin.Context) {
	response, count, err := s.WinnerService.GetAllWinners(0, 0)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Winners []map[string]interface{}
	for _, winner := range response {
		Winners = append(Winners, dto.CreateWinnerResponse(&winner))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"winners": Winners, "count": count})
}
