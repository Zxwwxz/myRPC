package limit

type LimitInterface interface {
	Allow() bool
}
