package ipx

// AddrError declares the address error
type AddrError struct {
	Err  string
	Addr string
}

// Error throws the error
func (e *AddrError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Err
	if e.Addr != "" {
		s = "address " + e.Addr + ": " + s
	}
	return s
}
