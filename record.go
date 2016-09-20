package router

// record describes information about the way in which there are parameters.
type record struct {
	params  uint16      // the number of parameters
	parts   []string    // way disassembled into its component parts
	handler interface{} // the request handler or something that is connected with it
}

// records describes a list of parameters and supports sorting on the number of
// parameters: the smaller the parameters the higher in the list. Account
// dynamic parameter with the lowest priority, i.e. places them in the end of
// the list.
type records []*record

// support methods for sorting.
func (n records) Len() int           { return len(n) }
func (n records) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n records) Less(i, j int) bool { return n[i].params < n[j].params }
