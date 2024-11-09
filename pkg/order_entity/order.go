package order_entity

import (
	"database/sql"
	"time"
)

type Order struct {
	ID               int64        `json:"ID" db:"id"`
	UserID           int64        `json:"userID" db:"user_id"`
	OrderCreateDate  time.Time    `json:"orderCreateDate" db:"created_at"`
	ExpiryDate       time.Time    `json:"expiryDate" db:"expiry_date"`
	AcceptDate       sql.NullTime `json:"acceptDate" db:"accept_return_order_date"`
	ReturnFromClient sql.NullTime `json:"returnFromClient" db:"returned_from_client_at"`
	ReturnToCourier  sql.NullTime `json:"returnToCourier" db:"returned_to_courier_at"`
	Packaging        string       `json:"packaging" db:"packaging"`
	Weigh            float32      `json:"weigh" db:"weigh"`
	Cost             float32      `json:"cost" db:"cost"`
}
