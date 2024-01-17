package dto

import "sorteio-api/source/models"

type SorteioCreateRequestDto struct {
	Name          string `form:"name" json:"name" xml:"name" binding:"required"`
	Description   string `form:"description" json:"description" xml:"description" binding:"required"`
	TicketsAmount string `form:"ticketsAmount" json:"ticketsAmount" xml:"ticketsAmount" binding:"required"`
	TicketsPrice  string `form:"ticketsPrice" json:"ticketsPrice" xml:"ticketsPrice" binding:"required"`
	//	TicketsStart             string `form:"ticketsStart" json:"ticketsStart" xml:"ticketsStart" binding:"required"`
	FinishDate               string `form:"finishDate" json:"finishDate" xml:"finishDate" binding:"required"`
	DisplayFinishDate        string `form:"display_finish_date" json:"display_finish_date" xml:"display_finish_date" binding:"required"`
	MinimalTicketForOrder    string `form:"minimalTicketForOrder" json:"minimalTicketForOrder" xml:"minimalTicketForOrder" binding:"required"`
	MaximumTicketForOrder    string `form:"maximumTicketForOrder" json:"maximumTicketForOrder" xml:"maximumTicketForOrder" binding:"required"`
	LimitTimeToExpireOrder   string `form:"limitTimeToExpireOrder" json:"limitTimeToExpireOrder" xml:"limitTimeToExpireOrder"`
	TicketsMinimalToDiscount string `form:"ticketsMinimalToDiscount" json:"ticketsMinimalToDiscount" xml:"ticketsMinimalToDiscount" binding:"required"`
	DiscountAmount           string `form:"discountAmount" json:"discountAmount" xml:"discountAmount" binding:"required"`
	Status                   string `form:"status" json:"status" xml:"status" binding:"required"`
}

func CreateSorteioResponse(sorteio *models.Sorteio) map[string]interface{} {

	handleData := map[string]interface{}{
		"id":          sorteio.SorteioId,
		"name":        sorteio.Name,
		"images":      sorteio.Images,
		"description": sorteio.Description,
		"tickets": map[string]interface{}{
			"amount":          sorteio.Tickets.Amount,
			"price":           sorteio.Tickets.Price,
			"minimalForOrder": sorteio.Tickets.MinimalForOrder,
			"maximumForOrder": sorteio.Tickets.MaximumForOrder,
		},
		"discount": map[string]interface{}{
			"amount":  sorteio.Discount.Amount,
			"minimal": sorteio.Discount.Minimal,
		},
		"percentage":          sorteio.Percentage,
		"earning":             sorteio.Earning,
		"status":              sorteio.Status,
		"createdAt":           sorteio.CreatedAt,
		"finishDate":          sorteio.FinishDate,
		"ticketsSold":         sorteio.TicketsSold,
		"display_finish_date": sorteio.DisplayFinishDate,
	}

	return handleData
}

type UpdateSorteioDto struct {
	SorteioId     string `form:"id" json:"id" xml:"id" binding:"required"`
	Name          string `form:"name" json:"name" xml:"name" binding:"required"`
	Description   string `form:"description" json:"description" xml:"description" binding:"required"`
	TicketsAmount string `form:"ticketsAmount" json:"ticketsAmount" xml:"ticketsAmount" binding:"required"`
	TicketsPrice  string `form:"ticketsPrice" json:"ticketsPrice" xml:"ticketsPrice" binding:"required"`
	//	TicketsStart             string `form:"ticketsStart" json:"ticketsStart" xml:"ticketsStart" binding:"required"`
	FinishDate               string `form:"finishDate" json:"finishDate" xml:"finishDate" binding:"required"`
	DisplayFinishDate        string `form:"display_finish_date" json:"display_finish_date" xml:"display_finish_date" binding:"required"`
	MinimalTicketForOrder    string `form:"minimalTicketForOrder" json:"minimalTicketForOrder" xml:"minimalTicketForOrder" binding:"required"`
	MaximumTicketForOrder    string `form:"maximumTicketForOrder" json:"maximumTicketForOrder" xml:"maximumTicketForOrder" binding:"required"`
	LimitTimeToExpireOrder   string `form:"limitTimeToExpireOrder" json:"limitTimeToExpireOrder" xml:"limitTimeToExpireOrder" `
	TicketsMinimalToDiscount string `form:"ticketsMinimalToDiscount" json:"ticketsMinimalToDiscount" xml:"ticketsMinimalToDiscount" binding:"required"`
	DiscountAmount           string `form:"discountAmount" json:"discountAmount" xml:"discountAmount" binding:"required"`
	Status                   string `form:"status" json:"status" xml:"status" binding:"required"`
}
