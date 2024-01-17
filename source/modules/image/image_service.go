package image

import (
	"sorteio-api/source/models"
	"sorteio-api/source/utils"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
)

// SERVICE
type IImageService interface {
	CreateImage(pathUrl string) interface{}
	GetAImage(imageId string) string
}

type ImageService struct {
	ImageRepository IImageRepository
	Validate        *validator.Validate
}

func InitImageService(ImageRepository IImageRepository, validate *validator.Validate) IImageService {
	return &ImageService{
		ImageRepository: ImageRepository,
		Validate:        validate,
	}
}

func (image *ImageService) CreateImage(pathUrl string) interface{} {
	var imageSingle = models.Image{
		PathUrl: pathUrl,
	}

	if validationErr := image.Validate.Struct(&imageSingle); validationErr != nil {
		sentry.CaptureException(validationErr)
		panic(validationErr)
	}

	var resultImage, _ = image.ImageRepository.Create(imageSingle)

	return resultImage
}

func (image *ImageService) GetAImage(imageId string) string {
	var ImageResult = image.ImageRepository.FindOne(imageId)

	return utils.NormalizeImageFromDB(ImageResult.PathUrl)
}
