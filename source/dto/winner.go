package dto

import "sorteio-api/source/models"

type WinnerCreateResponseDto struct {
	Name        string `form:"name" json:"name" xml:"name" binding:"required"`
	UserId      string `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	SorteioId   string `form:"sorteio_id" json:"sorteio_id" xml:"sorteio_id" binding:"required"`
	SorteioName string `form:"sorteio_name" json:"sorteio_name" xml:"sorteio_name" binding:"required"`
	CotaNumber  string `form:"cota_number" json:"cota_number" xml:"cota_number" binding:"required"`
	Image       string `form:"image" json:"image" xml:"image" binding:"required"`
}

func CreateWinnerResponse(sorteio *models.Winner) map[string]interface{} {

	handleData := map[string]interface{}{
		"id":          sorteio.ID,
		"name":        sorteio.Name,
		"userId":      sorteio.UserId,
		"sorteioId":   sorteio.SorteioId,
		"sorteioName": sorteio.SorteioName,
		"cotaNumber":  sorteio.CotaNumber,
		"image":       sorteio.Image,
	}

	return handleData
}

type WinnerCreateRequestDto struct {
	Name        string `form:"name" json:"name" xml:"name" binding:"required"`
	UserId      string `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	SorteioId   string `form:"sorteio_id" json:"sorteio_id" xml:"sorteio_id" binding:"required"`
	SorteioName string `form:"sorteio_name" json:"sorteio_name" xml:"sorteio_name" binding:"required"`
	CotaNumber  string `form:"cota_number" json:"cota_number" xml:"cota_number" binding:"required"`
	Image       string `form:"image" json:"image" xml:"image" binding:"required"`
}
