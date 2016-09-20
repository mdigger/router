# A simple router with support for named parameters

[![GoDoc](https://godoc.org/github.com/mdigger/router?status.svg)](https://godoc.org/github.com/mdigger/router)
[![Build Status](https://travis-ci.org/mdigger/router.svg)](https://travis-ci.org/mdigger/router)
[![Coverage Status](https://coveralls.io/repos/github/mdigger/router/badge.svg?branch=master)](https://coveralls.io/github/mdigger/router?branch=master)

router contains a fairly simple generic "router", which is the workpiece to replace `http.ServeMux` with support for named parameters in the path.

The current implementation allows you to associate with named routes any objects not only handlers of HTTP requests. For named parameters the marker is used `:`, dynamic — `*`. Though, it's all custom. Of course, static path without any parameters too supported.

Examples of questions ways:

	/user/:name
	/user/test
	/files/*filename
	/repos/:owner/:repo/pulls

### How to use?

The library is rather abstract and is not tied directly to `http`. So it is possible 
to use not only for its direct purpose:

```go
package main

import (
	"fmt"
	"github.com/mdigger/router"
)

func main() {
	var paths router.Paths
	paths.Add("/users", "usersList")
	paths.Add("/users/:name", "userName")
	paths.Add("/users/me", "userMe")

	name, params := paths.Lookup("/users/mdigger")
	fmt.Printf("%v: %v\n", name, params.Get("name"))
	// Output: userName: mdigger
}
```

### Excuse from any warranty

Previously, it was an integral part of the library <github.com/mdigger/rest> where, for the most part, all of this functionality was just hidden and not available for self-use. But it took me for some internal projects and I decided to submit it in a separate library. I don't guarantee that the library will from time to time to change my their own needs, so if you want to use it in their projects the best way is to take it entirely and continue to do with it everything that you want.
