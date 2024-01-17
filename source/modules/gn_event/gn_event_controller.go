package gnEvent

// CONTROLLER
type IGnEventController interface {
}

type GnEventController struct {
	GnEventService IGnEventService
}

func InitGnEventController(GnEventService IGnEventService) IGnEventController {
	return &GnEventController{
		GnEventService: GnEventService,
	}
}
