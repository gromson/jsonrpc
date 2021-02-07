package jsonrpc

// Namespace interface representing a namespace
type Namespace interface {
	// RegisterNS registers a namespace in a service
	RegisterNS(name string, ns Namespace)
	// Namespace retrieves a namespace from a service
	Namespace(name string) (Namespace, Error)
	// Register register a method in the namespace
	Register(name string, callback Callback)
	// Callback register a method in the namespace
	Callback(name string) (Callback, Error)
}

// Callback implementation of the RPC method
type Callback func(message []byte) (interface{}, Error)

type namespace struct {
	namespaces map[string]Namespace
	callbacks  map[string]Callback
}

// NewNamespace creates namespace
func NewNamespace() Namespace {
	namespaces := make(map[string]Namespace, 2)
	callbacks := make(map[string]Callback, 2)

	return &namespace{namespaces: namespaces, callbacks: callbacks}
}

// RegisterNS registers a namespace in a service
func (ns *namespace) RegisterNS(name string, sns Namespace) {
	ns.namespaces[name] = sns
}

// Namespace retrieves a namespace from a service
func (ns *namespace) Namespace(name string) (Namespace, Error) {
	item, ok := ns.namespaces[name]

	if !ok {
		return nil, NewNsNotFoundError(name, nil)
	}

	return item, nil
}

// Register register a method in the namespace
func (ns *namespace) Register(name string, callback Callback) {
	ns.callbacks[name] = callback
}

// Callback retrieves a method from the namespace
func (ns *namespace) Callback(name string) (Callback, Error) {
	fn, ok := ns.callbacks[name]

	if !ok {
		return nil, NewMethodNotFoundError(name, nil)
	}

	return fn, nil
}
