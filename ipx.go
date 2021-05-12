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

// A ParseError is the error type of literal network address parsers.
type ParseError struct {
	// Type is the type of string that was expected, such as
	// "IP address", "CIDR address".
	Type string

	// Text is the malformed text string.
	Text string
}

// Error ..
func (e *ParseError) Error() string { return "invalid " + e.Type + ": " + e.Text }

// Timeout ..
func (e *ParseError) Timeout() bool { return false }

// Temporary ..
func (e *ParseError) Temporary() bool { return false }
