package apierrors

type Error int

const (
	BadRequest    Error = 4001
	OrderNotFound Error = 4004
	Panic         Error = 9009
)
