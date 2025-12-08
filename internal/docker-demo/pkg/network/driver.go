package network

type Driver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network *Network) error
	Connect(networkName string, endpoint *Endpoint) error
	Disconnect(endpointID string) error
}
