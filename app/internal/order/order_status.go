package order

type Status string

const (
	Placed     Status = "order_placed"
	InProgress Status = "order_in_progress"
	Canceled   Status = "order_canceled"
	Completed  Status = "order_completed"
)
