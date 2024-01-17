package winner

import (
	"context"
	DBConnection "sorteio-api/source/database"
	"sorteio-api/source/dto"
	"sorteio-api/source/models"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SERVICE
type IWinnerService interface {
	GetAllWinners(page int, pageSize int) ([]models.Winner, int, error)
	CreateWinner(data dto.WinnerCreateRequestDto) (interface{}, error)
	GetAWinnerFromId(winnerId string) (models.Winner, bool)
}

type WinnerService struct {
	WinnerRepository IWinnerRepository
	Validate         *validator.Validate
}

func InitWinnerService(WinnerRepository IWinnerRepository, validate *validator.Validate) IWinnerService {
	return &WinnerService{
		WinnerRepository: WinnerRepository,
		Validate:         validate,
	}
}

func (s *WinnerService) GetAllWinners(page int, pageSize int) ([]models.Winner, int, error) {
	filter := bson.M{}

	results, err, ctx := DBConnection.FindMultipleDocuments("winner", filter, page, pageSize)

	count := results.RemainingBatchLength()

	var Winners []models.Winner

	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleWinner models.Winner
		if err = results.Decode(&singleWinner); err != nil {
			panic(err)
		}

		Winners = append(Winners, singleWinner)
	}

	return Winners, count, err
}

func (s *WinnerService) CreateWinner(data dto.WinnerCreateRequestDto) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	coll := DBConnection.GetCollection("winner")
	defer cancel()

	if validationErr := s.Validate.Struct(&data); validationErr != nil {
		panic(validationErr)
	}

	userId, _ := primitive.ObjectIDFromHex(data.UserId)

	var Winner = models.Winner{
		ID:          primitive.NewObjectID(),
		Name:        data.Name,
		UserId:      userId,
		SorteioId:   data.SorteioId,
		SorteioName: data.SorteioName,
		CotaNumber:  data.CotaNumber,
		Image:       data.Image,
		CreatedAt:   time.Now(),
	}

	result, err := coll.InsertOne(ctx, Winner)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return nil, err
	}

	return result.InsertedID, err
}

func (s *WinnerService) GetAWinnerFromId(winnerId string) (models.Winner, bool) {
	var winner models.Winner

	objId, _ := primitive.ObjectIDFromHex(winnerId)

	result, _ := DBConnection.FindADocument("winner", bson.M{"_id": objId})
	err := result.Decode(&winner)

	if err != nil {
		return models.Winner{}, false
	}

	return winner, true
}
