package router_test

import (
	"fmt"

	"github.com/mdigger/router"
)

func Example() {
	var paths router.Paths
	paths.Add("/users", "usersList")
	paths.Add("/users/:name", "userName")
	paths.Add("/users/me", "userMe")

	name, params := paths.Lookup("/users/mdigger")
	fmt.Printf("%v: %v\n", name, params.Get("name"))
	// Output: userName: mdigger
}
