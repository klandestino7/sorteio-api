package ticket

import (
	"net/http"
	"sorteio-api/source/modules/session"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CONTROLLER
type ITicketController interface {
	TryGetTicket(c *gin.Context)
}

type TicketController struct {
	TicketService  ITicketService
	SessionService session.ISessionService
}

func InitTicketController(TicketService ITicketService, SessionService session.ISessionService) ITicketController {
	return &TicketController{
		TicketService:  TicketService,
		SessionService: SessionService,
	}
}

func (tc *TicketController) TryGetTicket(c *gin.Context) {
	ticketId := c.Query("ticketId")
	sorteioId := c.Query("sorteioId")

	num, _ := strconv.Atoi(ticketId)
	sortId, _ := strconv.ParseInt(sorteioId, 10, 64)

	var sum int = 0

	for sum < num {
		ticketNumber := tc.TicketService.GenerateTicketNumber(sortId, num)

		if ticketNumber != -1 {
			sum += 1
		}

	}

	c.IndentedJSON(http.StatusOK, gin.H{"response": "ok"})
}
