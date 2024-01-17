package sorteio

import (
	"fmt"
	"net/http"
	"sorteio-api/source/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CONTROLLER
type ISorteioController interface {
	RequestSorteioFromId(c *gin.Context)
	RequestAllSorteios(c *gin.Context)
	RequestFullSorteios(c *gin.Context)
	RequestCreateNewSorteio(c *gin.Context)
	RequestGenerateWinner(c *gin.Context)
	RequestUpdateSorteio(c *gin.Context)
}

type SorteioController struct {
	SorteioService ISorteioService
}

func InitSorteioController(SorteioService ISorteioService) ISorteioController {
	return &SorteioController{
		SorteioService: SorteioService,
	}
}

func (ct *SorteioController) RequestSorteioFromId(c *gin.Context) {
	sorteioId := c.Param("sorteioId")

	if sorteioId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorteioId Invalid"})
		return
	}

	response, status := ct.SorteioService.GetASorteioFromId(sorteioId)

	if !status {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorteio Invalid"})
		return
	}

	sorteio := dto.CreateSorteioResponse(&response)

	c.IndentedJSON(http.StatusOK, gin.H{"sorteio": sorteio})
}

func (ct *SorteioController) RequestAllSorteios(c *gin.Context) {
	response, count, err := ct.SorteioService.GetAllSorteios(0, 0)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Sorteios []map[string]interface{}
	for _, sorteio := range response {
		Sorteios = append(Sorteios, dto.CreateSorteioResponse(&sorteio))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"sorteios": Sorteios, "count": count})
}

func (ct *SorteioController) RequestFullSorteios(c *gin.Context) {
	response, count, err := ct.SorteioService.GetFullSorteios(0, 0)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		return
	}

	var Sorteios []map[string]interface{}
	for _, sorteio := range response {
		Sorteios = append(Sorteios, dto.CreateSorteioResponse(&sorteio))
	}

	c.IndentedJSON(http.StatusOK, gin.H{"sorteios": Sorteios, "count": count})
}

func (ct *SorteioController) RequestCreateNewSorteio(c *gin.Context) {
	var json dto.SorteioCreateRequestDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		panic(err)
		return
	}

	result, _ := ct.SorteioService.CreateSorteio(json)

	if primitive.IsValidObjectID(fmt.Sprintf("%s", result)) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func (ct *SorteioController) RequestGenerateWinner(c *gin.Context) {
	sorteioId := c.Param("sorteioId")

	if sorteioId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorteioId Invalid"})
		return
	}

	response := ct.SorteioService.GenerateAWinnerToSorteio(sorteioId)

	c.IndentedJSON(http.StatusOK, gin.H{"ticket": dto.CreateTicketWinnerGenerateResponse(response)})
}

func (ct *SorteioController) RequestUpdateSorteio(c *gin.Context) {
	var json dto.UpdateSorteioDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateBadRequestErrorDto(err))
		panic(err)
		return
	}

	ct.SorteioService.UpdateSorteio(json)

	// if primitive.IsValidObjectID(fmt.Sprintf("%s", result)) {
	c.JSON(http.StatusOK, gin.H{"success": true})
	// }
}
