package limiter

type LimitInterface interface {
	Allow() bool
}
