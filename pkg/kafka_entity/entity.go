package kafka_entity

type AcceptReturnEventMessage struct {
	OrderID int64 `json:"order_id"`
	UserID  int64 `json:"user_id"`
}
type AddOrderEventMessage struct {
	OrderID   int64   `json:"order_id"`
	UserID    int64   `json:"user_id"`
	ShelfLife int64   `json:"shelf_life"`
	Packaging string  `json:"packaging"`
	Weigh     float32 `json:"weigh"`
	Cost      float32 `json:"cost"`
}
type ReturnOrderEventMessage struct {
	OrderID int64 `json:"order_id"`
}
type UpdateIssuedEventMessage struct {
	OrderIDs []string `json:"order_id"`
}
