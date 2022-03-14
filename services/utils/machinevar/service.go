package machinevar

type Service interface {

	// get environment variable
	// will result error if not found or empty
	GetVar(key string) (result string, err error)
}
