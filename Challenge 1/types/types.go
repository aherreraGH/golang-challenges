package types

/**
 * package delivery routes
 */
type Package struct {
	ID       int
	Address  string
	Distance int
}

/**
 * Worker pooling results
 */
type DeliveryResult struct {
	PackageID int
	WorkerID  int
	Status    string
}
