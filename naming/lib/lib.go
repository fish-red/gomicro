package lib

import (
	"google.golang.org/grpc/naming"
)

// GenUpdates generates the added and deleted updates
func GenUpdates(a, b []string) []*naming.Update {
	updates := []*naming.Update{}

	// delete the elements in a and not in b
	deleted := diff(a, b)
	for _, addr := range deleted {
		update := &naming.Update{Op: naming.Delete, Addr: addr}
		updates = append(updates, update)
	}

	// add the elements in b and not in a
	added := diff(b, a)
	for _, addr := range added {
		update := &naming.Update{Op: naming.Add, Addr: addr}
		updates = append(updates, update)
	}
	return updates
}

// diff(a, b) = a - a(n)b
func diff(a, b []string) []string {
	var d []string
	for _, va := range a {
		found := false
		for _, vb := range b {
			if va == vb {
				found = true
				break
			}
		}
		if !found {
			d = append(d, va)
		}
	}
	return d
}
