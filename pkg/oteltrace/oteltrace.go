package oteltrace

type Config interface {
	config()
}

type Option func(Config)
