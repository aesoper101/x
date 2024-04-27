package configx

type Config interface {
	Validate() error
}
