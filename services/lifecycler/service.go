package lifecycler

type Service interface {
	// Executed before builder options
	PreConfig() (err error)
	// Executed after builder options
	PostConfig() (err error)
	// Function to start the lifecycler
	Run() (err error)
	AddProviders(providers ...interface{}) (err error)
	AddInvokers(invokers ...interface{}) (err error)
}
