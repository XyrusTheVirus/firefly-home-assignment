package transporters

type TransporterInterface interface {
	// Transport sends the processed data to the desired destination and returns the response and error if any.
	// It sustains all ing of protocols responses, hence, defines as interface{} and the caller should type assert it.
	Transport() (interface{}, error)
}
