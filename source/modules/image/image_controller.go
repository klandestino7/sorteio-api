package image

// CONTROLLER
type IImageController interface {
}

type ImageController struct {
	ImageService IImageService
}

func InitImageController(ImageService IImageService) IImageController {
	return &ImageController{
		ImageService: ImageService,
	}
}
