package gnEvent

import (
	"context"
	"encoding/json"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/dto"
	"sorteio-api/source/models"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SERVICE
type IGnEventService interface {
	CreateGnEvent(dataHandle dto.PixResponse) (interface{}, error)
}

type GnEventService struct {
	GnEventRepository IGnEventRepository
	Validate          *validator.Validate
}

func InitGnEventService(GnEventRepository IGnEventRepository, validate *validator.Validate) IGnEventService {
	return &GnEventService{
		GnEventRepository: GnEventRepository,
		Validate:          validate,
	}
}

func (s *GnEventService) CreateGnEvent(dataHandle dto.PixResponse) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("gnevent")
	defer cancel()

	b, err := json.Marshal(dataHandle)

	var gnEvent = models.GnEvent{
		ID:    primitive.NewObjectID(),
		Eid:   dataHandle.EndToEndId,
		Txid:  dataHandle.Txid,
		Event: string(b),
	}

	if validationErr := s.Validate.Struct(&dataHandle); validationErr != nil {
		panic(validationErr)
	}

	result, err := coll.InsertOne(ctx, gnEvent)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, err
	}

	return result.InsertedID, err
}
