package Jobs

type Jobs interface {
	Execute() Error
}