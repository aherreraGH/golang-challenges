package types

type Package struct {
	ID       int
	Address  string
	Distance int
}

type DeliveryResult struct {
	PackageID int
	WorkerID  int
	Status    string
}
