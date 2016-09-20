// Package router contains a fairly simple generic "router", which is the
// workpiece for replacement http.ServeMux with support for named parameters in
// the path.
//
// The current implementation allows you to associate with named routes any
// objects not only handlers of HTTP requests. For named parameters used the
// marker `:`, dynamic — `*`. Though, it's all custom. Of course, static path
// without any parameters too supported.
//
// Examples of questions ways:
// 	/user/:name
//	/user/test
// 	/files/*filename
// 	/repos/:owner/:repo/pulls
//
// Excuse from any warranty
//
// Previously, it was an integral part of the library github.com/mdigger/rest
// where, for the most part, all of this functionality was just hidden and not
// available for self-use. But it took me for some internal projects and I
// decided to submit it in a separate library. I don't guarantee that the
// library will from time to time to change my their own needs, so if you want
// to use it in their projects the best way is to take it entirely and continue
// to do with it everything that you want.
package router
