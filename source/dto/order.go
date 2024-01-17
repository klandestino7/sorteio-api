package dto

import "sorteio-api/source/models"

type OrderRequestDto struct {
	Name              string `form:"name" json:"name" xml:"name" binding:"required"`
	Email             string `form:"email" json:"email" xml:"email" binding:"required"`
	Phone             string `form:"phone" json:"phone" xml:"phone" binding:"required"`
	PhoneConfirmation string `form:"phone_confirmation" json:"phone_confirmation" xml:"phone_confirmation" binding:"required"`
	Cpf               string `form:"cpf" json:"cpf" xml:"cpf" binding:"required"`
	Tickets           int    `form:"tickets" json:"tickets" xml:"tickets" binding:"required"`
	Sorteio           int    `form:"sorteio" json:"sorteio" xml:"sorteio" binding:"required"`
	Referal           string `form:"referal" json:"referal" xml:"referal"`
}

type OrderResponseDto struct {
	OrderId       int64  `bson:"number"`
	SorteioId     string `bson:"sorteio_id"`
	Tickets       []int  `bson:"tickets"`
	TicketsAmount int    `bson:"tickets_amount"`
	TransaciontId string `bson:"transaction_id"`
	Status        string `bson:"status"`
	Total         int    `total`
	ExpireAt      int64  `expire_at`
}

type OrderRequestFromCredential struct {
	Credential string `form:"credential" json:"credential" xml:"credential" binding:"required"`
	Value      string `form:"value" json:"value" xml:"value" binding:"required"`
}

func CreateOrderResponse(order *models.Order) map[string]interface{} {

	handleData := map[string]interface{}{
		"id":            order.ID,
		"orderId":       order.OrderId,
		"userId":        order.UserId,
		"userName":      order.UserName,
		"sorteioId":     order.SorteioId,
		"transactionId": order.TransactionId,
		"tickets":       order.Tickets,
		"ticketsAmount": order.TicketsAmount,
		"qrCode":        order.QRCodeString,
		"qrCodeImage":   order.QRCodeImage,
		"status":        order.Status,
		"paymentMethod": order.PaymentMethod,
		"total":         order.Total,
		"referal":       order.Referal,
		"createdAt":     order.CreatedAt,
		"expireAt":      order.ExpireAt,
	}

	return handleData
}

type PixHandleDto struct {
	Parametros map[string]interface{} `form:"parametros" json:"parametros" xml:"parametros"`
	Pix        []PixResponse          `form:"pix" json:"pix" xml:"pix" binding:"required"`
}

type PixResponse struct {
	EndToEndId string         `json:"endToEndId" binding:"required"`
	Txid       string         `json:"txid" binding:"required"`
	Valor      string         `json:"valor" binding:"required" `
	Chave      string         `json:"chave" binding:"required" `
	Horario    string         `json:"horario" binding:"required"`
	Devolucoes []PixDevolucao `json:"devolucoes"`
}

type PixDevolucao struct {
	Id      string                 `json:"id" binding:"required"`
	RtrId   string                 `json:"rtrId" binding:"required"`
	Valor   string                 `json:"valor" binding:"required"`
	Horario map[string]interface{} `json:"horario" binding:"required"`
	Status  string                 `json:"status" binding:"required"`
}

func CreateTicketWinnerGenerateResponse(ticket models.Ticket) map[string]interface{} {
	return map[string]interface{}{
		"id":             ticket.ID,
		"number":         ticket.Number,
		"user_id":        ticket.UserId,
		"sorteio_id":     ticket.SorteioId,
		"order_id":       ticket.OrderId,
		"status":         ticket.Status,
		"created_at":     ticket.CreatedAt,
		"transaction_id": ticket.TransactionId,
	}
}
