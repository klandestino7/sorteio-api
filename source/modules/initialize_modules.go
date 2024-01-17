package modules

import (
	authToken "sorteio-api/source/modules/auth_token"
	gnEvent "sorteio-api/source/modules/gn_event"
	"sorteio-api/source/modules/image"
	"sorteio-api/source/modules/order"
	"sorteio-api/source/modules/session"
	"sorteio-api/source/modules/sorteio"
	"sorteio-api/source/modules/ticket"
	"sorteio-api/source/modules/user"
	"sorteio-api/source/modules/winner"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var AuthTokenRepository = authToken.InitAuthTokenRepository()
var AuthTokenService = authToken.InitAuthTokenService(AuthTokenRepository, validate)
var AuthTokenController = authToken.InitAuthTokenController(AuthTokenService)

var GnEventRepository = gnEvent.InitGnEventRepository()
var GnEventService = gnEvent.InitGnEventService(GnEventRepository, validate)
var GnEventController = gnEvent.InitGnEventController(GnEventService)

var ImageRepository = image.InitImageRepository()
var ImageService = image.InitImageService(ImageRepository, validate)
var ImageControler = image.InitImageController(ImageService)

var WinnerRepositoy = winner.InitWinnerRepository()
var WinnerService = winner.InitWinnerService(WinnerRepositoy, validate)
var WinnerController = winner.InitWinnerController(WinnerService)

var UserRepository = user.InitUserRepository()
var SessionRepository = session.InitSessionRepository()
var UserService = user.InitUserService(UserRepository, validate)
var SessionService = session.InitSessionService(SessionRepository, validate)
var UserController = user.InitUserController(UserService, SessionService)
var SessionController = session.InitSessionController(SessionService)

var TicketRepository = ticket.InitTicketRepository()
var TicketService = ticket.InitTicketService(TicketRepository, validate)
var TicketController = ticket.InitTicketController(TicketService, SessionService)

var SorteioRepository = sorteio.InitSorteioRepository()
var SorteioService = sorteio.InitSorteioService(SorteioRepository, validate, TicketService, WinnerService)
var SorteioController = sorteio.InitSorteioController(SorteioService)

var OrderRepository = order.InitOrderRepository()
var OrderService = order.InitOrderService(OrderRepository, validate, UserService, SorteioService, TicketService)
var OrderController = order.InitOrderController(OrderService, GnEventService)
