package ipx

// IPNet represents an IP network.
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}

// Contains reports whether the network includes ip.
func (n *IPNet) Contains(ip IP) bool {

}
