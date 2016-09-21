package router

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	testRoutersWithPaths(t, []string{
		"/",
		"/path/to/route",
		"/path/to/other",
		"/path/to/route/a",
		"/path/to/:param",
		"/path/to/wildcard/*routepath",
		"/path/to/:param1/:param2",
		"/path/to/:param1/sep/:param2",
		"/:year/:month/:day",
		"/user/:id",
		"/a/to/b/:param/*routepath",
	}, []TestStruct{
		{"/", 0, nil},
		{"/path/to/route", 1, nil},
		{"/path/to/other", 2, nil},
		{"/path/to/route/a", 3, nil},
		{"/path/to/hoge", 4, Params{{"param", "hoge"}}},
		{"/path/to/wildcard/some/params", 5, Params{{"routepath", "some/params"}}},
		{"/path/to/o1/o2", 6, Params{{"param1", "o1"}, {"param2", "o2"}}},
		{"/path/to/p1/sep/p2", 7, Params{{"param1", "p1"}, {"param2", "p2"}}},
		{"/2014/01/06", 8, Params{{"year", "2014"}, {"month", "01"}, {"day", "06"}}},
		{"/user/777", 9, Params{{"id", "777"}}},
		{"/a/to/b/p1/some/wildcard/params", 10, Params{{"param", "p1"}, {"routepath", "some/wildcard/params"}}},
		{"/missing", nil, nil},
	})

	testRoutersWithPaths(t, []string{
		"/",
		"/:b",
		"/*wildcard",
	}, []TestStruct{
		{"/", 0, nil},
		{"/true", 1, Params{{"b", "true"}}},
		{"/foo/bar", 2, Params{{"wildcard", "foo/bar"}}},
	})

	testRoutersWithPaths(t, []string{
		"/networks/:owner/:repo/events",
		"/orgs/:org/events",
		"/notifications/threads/:id",
	}, []TestStruct{
		{"/networks/:owner/:repo/events", 0, Params{{"owner", ":owner"}, {"repo", ":repo"}}},
		{"/orgs/:org/events", 1, Params{{"org", ":org"}}},
		{"/notifications/threads/:id", 2, Params{{"id", ":id"}}},
	})

	testRoutersWithPaths(t, []string{
		"/",
	}, []TestStruct{
		{"/user/alice", nil, nil},
	})

	testRoutersWithPaths(t, []string{
		"/user/:name",
	}, []TestStruct{
		{"/", nil, nil},
	})

	testRoutersWithPaths(t, []string{
		"/*wildcard",
		"/a/:b",
	}, []TestStruct{
		{"/a", 0, Params{{"wildcard", "a"}}},
	})

	testRoutersWithPaths(t, []string{
		"/:mx-name/mxproxy",
		"/:mx-name/store/*filename",
	}, []TestStruct{
		{"/xyzrd/mxproxy", 0, Params{{"mx-name", "xyzrd"}}},
		{"/xyzrd/mxproxy/", nil, nil},
		{"/xyzrd/test", nil, nil},
		{"/xyzrd/store", nil, nil},
		{"/xyzrd/store/file/name", 1, Params{{"mx-name", "xyzrd"}, {"filename", "file/name"}}},
		{"/xyzrd/store/file/name/", 1, Params{{"mx-name", "xyzrd"}, {"filename", "file/name/"}}},
	})

}

func TestMoreRouters(t *testing.T) {
	var tests = []struct {
		URL, ParamURL string
		Params
	}{
		{"/authorizations", "/authorizations", nil},
		{"/authorizations/1", "/authorizations/:id", Params{{"id", "1"}}},
		{"/applications/1/tokens/zohRoo7e", "/applications/:client_id/tokens/:access_token", Params{{"client_id", "1"}, {"access_token", "zohRoo7e"}}},
		{"/events", "/events", nil},
		{"/repos/mdigger/rest/events", "/repos/:owner/:repo/events", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/networks/mdigger/rest/events", "/networks/:owner/:repo/events", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/orgs/something/events", "/orgs/:org/events", Params{{"org", "something"}}},
		{"/users/mdigger/received_events", "/users/:user/received_events", Params{{"user", "mdigger"}}},
		{"/users/mdigger/received_events/public", "/users/:user/received_events/public", Params{{"user", "mdigger"}}},
		{"/users/mdigger/events", "/users/:user/events", Params{{"user", "mdigger"}}},
		{"/users/mdigger/events/public", "/users/:user/events/public", Params{{"user", "mdigger"}}},
		{"/users/mdigger/events/orgs/something", "/users/:user/events/orgs/:org", Params{{"user", "mdigger"}, {"org", "something"}}},
		{"/feeds", "/feeds", nil},
		{"/notifications", "/notifications", nil},
		{"/repos/mdigger/rest/notifications", "/repos/:owner/:repo/notifications", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/notifications/threads/1", "/notifications/threads/:id", Params{{"id", "1"}}},
		{"/notifications/threads/2/subscription", "/notifications/threads/:id/subscription", Params{{"id", "2"}}},
		{"/repos/mdigger/rest/stargazers", "/repos/:owner/:repo/stargazers", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/users/mdigger/starred", "/users/:user/starred", Params{{"user", "mdigger"}}},
		{"/user/starred", "/user/starred", nil},
		{"/user/starred/mdigger/rest", "/user/starred/:owner/:repo", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/subscribers", "/repos/:owner/:repo/subscribers", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/users/mdigger/subscriptions", "/users/:user/subscriptions", Params{{"user", "mdigger"}}},
		{"/user/subscriptions", "/user/subscriptions", nil},
		{"/repos/mdigger/rest/subscription", "/repos/:owner/:repo/subscription", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/user/subscriptions/mdigger/rest", "/user/subscriptions/:owner/:repo", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/users/mdigger/gists", "/users/:user/gists", Params{{"user", "mdigger"}}},
		{"/gists", "/gists", nil},
		{"/gists/1", "/gists/:id", Params{{"id", "1"}}},
		{"/gists/2/star", "/gists/:id/star", Params{{"id", "2"}}},
		{"/repos/mdigger/rest/git/blobs/d30039aa3284fb929fcc80003234dedf93412aec", "/repos/:owner/:repo/git/blobs/:sha", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/repos/mdigger/rest/git/commits/d30039aa3284fb929fcc80003234dedf93412aec", "/repos/:owner/:repo/git/commits/:sha", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/repos/mdigger/rest/git/refs", "/repos/:owner/:repo/git/refs", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/git/tags/d30039aa3284fb929fcc80003234dedf93412aec", "/repos/:owner/:repo/git/tags/:sha", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/repos/mdigger/rest/git/trees/d30039aa3284fb929fcc80003234dedf93412aec", "/repos/:owner/:repo/git/trees/:sha", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/issues", "/issues", nil},
		{"/user/issues", "/user/issues", nil},
		{"/orgs/something/issues", "/orgs/:org/issues", Params{{"org", "something"}}},
		{"/repos/mdigger/rest/issues", "/repos/:owner/:repo/issues", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/issues/1", "/repos/:owner/:repo/issues/:number", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/assignees", "/repos/:owner/:repo/assignees", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/assignees/foo", "/repos/:owner/:repo/assignees/:assignee", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"assignee", "foo"}}},
		{"/repos/mdigger/rest/issues/1/comments", "/repos/:owner/:repo/issues/:number/comments", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/issues/1/events", "/repos/:owner/:repo/issues/:number/events", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/labels", "/repos/:owner/:repo/labels", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/labels/bug", "/repos/:owner/:repo/labels/:name", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"name", "bug"}}},
		{"/repos/mdigger/rest/issues/1/labels", "/repos/:owner/:repo/issues/:number/labels", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/milestones/1/labels", "/repos/:owner/:repo/milestones/:number/labels", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/milestones", "/repos/:owner/:repo/milestones", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/milestones/1", "/repos/:owner/:repo/milestones/:number", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/emojis", "/emojis", nil},
		{"/gitignore/templates", "/gitignore/templates", nil},
		{"/gitignore/templates/Go", "/gitignore/templates/:name", Params{{"name", "Go"}}},
		{"/meta", "/meta", nil},
		{"/rate_limit", "/rate_limit", nil},
		{"/users/mdigger/orgs", "/users/:user/orgs", Params{{"user", "mdigger"}}},
		{"/user/orgs", "/user/orgs", nil},
		{"/orgs/something", "/orgs/:org", Params{{"org", "something"}}},
		{"/orgs/something/members", "/orgs/:org/members", Params{{"org", "something"}}},
		{"/orgs/something/members/mdigger", "/orgs/:org/members/:user", Params{{"org", "something"}, {"user", "mdigger"}}},
		{"/orgs/something/public_members", "/orgs/:org/public_members", Params{{"org", "something"}}},
		{"/orgs/something/public_members/mdigger", "/orgs/:org/public_members/:user", Params{{"org", "something"}, {"user", "mdigger"}}},
		{"/orgs/something/teams", "/orgs/:org/teams", Params{{"org", "something"}}},
		{"/teams/1", "/teams/:id", Params{{"id", "1"}}},
		{"/teams/2/members", "/teams/:id/members", Params{{"id", "2"}}},
		{"/teams/3/members/mdigger", "/teams/:id/members/:user", Params{{"id", "3"}, {"user", "mdigger"}}},
		{"/teams/4/repos", "/teams/:id/repos", Params{{"id", "4"}}},
		{"/teams/5/repos/mdigger/rest", "/teams/:id/repos/:owner/:repo", Params{{"id", "5"}, {"owner", "mdigger"}, {"repo", "rest"}}},
		{"/user/teams", "/user/teams", nil},
		{"/repos/mdigger/rest/pulls", "/repos/:owner/:repo/pulls", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/pulls/1", "/repos/:owner/:repo/pulls/:number", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/pulls/1/commits", "/repos/:owner/:repo/pulls/:number/commits", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/pulls/1/files", "/repos/:owner/:repo/pulls/:number/files", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/pulls/1/merge", "/repos/:owner/:repo/pulls/:number/merge", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/repos/mdigger/rest/pulls/1/comments", "/repos/:owner/:repo/pulls/:number/comments", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"number", "1"}}},
		{"/user/repos", "/user/repos", nil},
		{"/users/mdigger/repos", "/users/:user/repos", Params{{"user", "mdigger"}}},
		{"/orgs/something/repos", "/orgs/:org/repos", Params{{"org", "something"}}},
		{"/repositories", "/repositories", nil},
		{"/repos/mdigger/rest", "/repos/:owner/:repo", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/contributors", "/repos/:owner/:repo/contributors", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/languages", "/repos/:owner/:repo/languages", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/teams", "/repos/:owner/:repo/teams", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/tags", "/repos/:owner/:repo/tags", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/branches", "/repos/:owner/:repo/branches", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/branches/master", "/repos/:owner/:repo/branches/:branch", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"branch", "master"}}},
		{"/repos/mdigger/rest/collaborators", "/repos/:owner/:repo/collaborators", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/collaborators/something", "/repos/:owner/:repo/collaborators/:user", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"user", "something"}}},
		{"/repos/mdigger/rest/comments", "/repos/:owner/:repo/comments", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/commits/d30039aa3284fb929fcc80003234dedf93412aec/comments", "/repos/:owner/:repo/commits/:sha/comments", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/repos/mdigger/rest/comments/1", "/repos/:owner/:repo/comments/:id", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "1"}}},
		{"/repos/mdigger/rest/commits", "/repos/:owner/:repo/commits", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/commits/d30039aa3284fb929fcc80003234dedf93412aec", "/repos/:owner/:repo/commits/:sha", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"sha", "d30039aa3284fb929fcc80003234dedf93412aec"}}},
		{"/repos/mdigger/rest/readme", "/repos/:owner/:repo/readme", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/keys", "/repos/:owner/:repo/keys", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/keys/1", "/repos/:owner/:repo/keys/:id", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "1"}}},
		{"/repos/mdigger/rest/downloads", "/repos/:owner/:repo/downloads", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/downloads/2", "/repos/:owner/:repo/downloads/:id", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "2"}}},
		{"/repos/mdigger/rest/forks", "/repos/:owner/:repo/forks", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/hooks", "/repos/:owner/:repo/hooks", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/hooks/2", "/repos/:owner/:repo/hooks/:id", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "2"}}},
		{"/repos/mdigger/rest/releases", "/repos/:owner/:repo/releases", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/releases/1", "/repos/:owner/:repo/releases/:id", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "1"}}},
		{"/repos/mdigger/rest/releases/1/assets", "/repos/:owner/:repo/releases/:id/assets", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"id", "1"}}},
		{"/repos/mdigger/rest/stats/contributors", "/repos/:owner/:repo/stats/contributors", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/stats/commit_activity", "/repos/:owner/:repo/stats/commit_activity", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/stats/code_frequency", "/repos/:owner/:repo/stats/code_frequency", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/stats/participation", "/repos/:owner/:repo/stats/participation", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/stats/punch_card", "/repos/:owner/:repo/stats/punch_card", Params{{"owner", "mdigger"}, {"repo", "rest"}}},
		{"/repos/mdigger/rest/statuses/master", "/repos/:owner/:repo/statuses/:ref", Params{{"owner", "mdigger"}, {"repo", "rest"}, {"ref", "master"}}},
		{"/search/repositories", "/search/repositories", nil},
		{"/search/code", "/search/code", nil},
		{"/search/issues", "/search/issues", nil},
		{"/search/users", "/search/users", nil},
		{"/legacy/issues/search/mdigger/rest/closed/test", "/legacy/issues/search/:owner/:repository/:state/:keyword", Params{{"owner", "mdigger"}, {"repository", "rest"}, {"state", "closed"}, {"keyword", "test"}}},
		{"/legacy/repos/search/test", "/legacy/repos/search/:keyword", Params{{"keyword", "test"}}},
		{"/legacy/user/search/test", "/legacy/user/search/:keyword", Params{{"keyword", "test"}}},
		{"/legacy/user/email/mdigger@xyzrd.com", "/legacy/user/email/:email", Params{{"email", "mdigger@xyzrd.com"}}},
		{"/users/mdigger", "/users/:user", Params{{"user", "mdigger"}}},
		{"/user", "/user", nil},
		{"/users", "/users", nil},
		{"/user/emails", "/user/emails", nil},
		{"/users/mdigger/followers", "/users/:user/followers", Params{{"user", "mdigger"}}},
		{"/user/followers", "/user/followers", nil},
		{"/users/mdigger/following", "/users/:user/following", Params{{"user", "mdigger"}}},
		{"/user/following", "/user/following", nil},
		{"/user/following/mdigger", "/user/following/:user", Params{{"user", "mdigger"}}},
		{"/users/mdigger/following/target", "/users/:user/following/:target_user", Params{{"user", "mdigger"}, {"target_user", "target"}}},
		{"/users/mdigger/keys", "/users/:user/keys", Params{{"user", "mdigger"}}},
		{"/user/keys", "/user/keys", nil},
		{"/user/keys/1", "/user/keys/:id", Params{{"id", "1"}}},
		{"/people/me", "/people/:userId", Params{{"userId", "me"}}},
		{"/people", "/people", nil},
		{"/activities/foo/people/vault", "/activities/:activityId/people/:collection", Params{{"activityId", "foo"}, {"collection", "vault"}}},
		{"/people/me/people/vault", "/people/:userId/people/:collection", Params{{"userId", "me"}, {"collection", "vault"}}},
		{"/people/me/openIdConnect", "/people/:userId/openIdConnect", Params{{"userId", "me"}}},
		{"/people/me/activities/vault", "/people/:userId/activities/:collection", Params{{"userId", "me"}, {"collection", "vault"}}},
		{"/activities/foo", "/activities/:activityId", Params{{"activityId", "foo"}}},
		{"/activities", "/activities", nil},
		{"/activities/foo/comments", "/activities/:activityId/comments", Params{{"activityId", "foo"}}},
		{"/comments/hoge", "/comments/:commentId", Params{{"commentId", "hoge"}}},
		{"/people/me/moments/vault", "/people/:userId/moments/:collection", Params{{"userId", "me"}, {"collection", "vault"}}},
	}

	urls := make([]string, len(tests))
	testStructs := make([]TestStruct, len(tests))
	for i, test := range tests {
		urls[i] = test.ParamURL
		testStructs[i] = TestStruct{
			URL:    test.URL,
			Index:  i,
			Params: test.Params,
		}
	}
	testRoutersWithPaths(t, urls, testStructs)
}

func TestMixPaths(t *testing.T) {
	tests := []string{
		"/user",
		"/user/:id",
		"/user/:id/post",
		"/user/:id/:group",
		"/user/:id/post/:cid",
		"/admin/:id/post/:cid",
		"/admin/:id/post/test",
		"/user/:id/post/:cid/:type",
	}
	var r Paths
	for i, url := range tests {
		err := r.Add(url, i)
		if err != nil {
			t.Error(err)
		}
	}
	// pretty.Println(r)
	for i := len(tests) - 1; i >= 0; i-- {
		index, params := r.Lookup(tests[i])
		if index != i {
			fmt.Println(i, index, tests[i], params)
			t.Errorf("не совпало с ожиданием: %v", i)
		}
		for _, param := range params {
			_ = param.String()
			if param.Value != params.Get(param.Key) {
				t.Errorf("bad param get by name: %v", param.Key)
			}
		}
		if params.Get("test") != "" {
			t.Error("bad param value `test`")
		}
	}

	// ---------------------
	// for cover
	// ---------------------
	if result := r.Path(len(tests)); result != nil {
		t.Errorf("bad path result: %v", result)
	}
	if r.Add("/test", nil) == nil {
		t.Error("add nil handler")
	}
	longPath := strings.Repeat("/test", 1<<15)
	if r.Add(longPath, 999) == nil {
		t.Error("add long path")
	}
	if r.Add("/file/*name/test", 998) == nil {
		t.Error("add bad catch all")
	}

}

func TestDoublePaths(t *testing.T) {
	tests := []string{
		// "/:user/:id/:id",
		// "/:user/:user/:id",
		"/:user/:name",
		"/:user/test",
	}
	var r Paths
	for i, url := range tests {
		err := r.Add(url, i)
		if err != nil {
			t.Error(err)
		}
	}
	// pretty.Println(r)
	for i, url := range tests {
		index, params := r.Lookup(url)
		if index != i {
			fmt.Println(i, index, tests[i], params)
			t.Errorf("не совпало с ожиданием: %v", i)
		}
	}
	index, params := r.Lookup("/user/name/vasya")
	if index != nil || len(params) != 0 {
		t.Error("bad lookup")
	}
}

func TestCatchAll(t *testing.T) {
	tests := []string{
		"/store/file/*filename",
		"/store/:file/*filename",
		"/store/file/test/*filename",
		"/store/:file/test/*filename",
	}
	var r Paths
	for i, url := range tests {
		err := r.Add(url, i)
		if err != nil {
			t.Error(err)
		}
	}
	// pretty.Println(r)
	for i := len(tests) - 1; i >= 0; i-- {
		index, params := r.Lookup(tests[i])
		if index != i {
			fmt.Println(i, index, tests[i], params)
			t.Errorf("не совпало с ожиданием: %v", i)
		}
	}
	for i, url := range []string{
		"/store/file/testfile",
		"/store/test/testfile",
		"/store/file/test/testfile",
		"/store/test/test/testfile",
	} {
		index, params := r.Lookup(url)
		if index != i {
			fmt.Println(i, index, tests[i], params)
			t.Errorf("не совпало с ожиданием: %v", i)
		}
	}
	index, params := r.Lookup("/store/test/test/testfile/2/1/3")
	if index != 3 {
		fmt.Println(index, tests[3], params)
		t.Errorf("не совпало с ожиданием: %v", 3)
	}
	index, params = r.Lookup("/store/file/test/testfile/2/1/3")
	if index != 2 {
		fmt.Println(index, tests[3], params)
		t.Errorf("не совпало с ожиданием: %v", 2)
	}

}

type TestStruct struct {
	URL    string      // адрес
	Index  interface{} // индекс совпадения
	Params             // найденные параметры и их значения
}

func testRoutersWithPaths(t *testing.T, routers []string, tests []TestStruct) {
	var r Paths
	for i, url := range routers {
		err := r.Add(url, i)
		if err != nil {
			t.Error(err)
		}
	}
	for i, url := range routers {
		paths := r.Path(i)
		path := "/" + strings.Join(paths, PathDelimeter)
		if url != path {
			t.Errorf("[%v] bad path: %v против %v", i, url, path)
		}
	}
	for i, obj := range tests {
		index, params := r.Lookup(obj.URL)
		if index != obj.Index {
			t.Errorf("[%v] не совпал выбор пути с ожиданием: %v против %v (%v)", i, index, obj.Index, obj.URL)
		}
		if !reflect.DeepEqual(params, obj.Params) {
			t.Errorf("[%v] не совпали параметры: %v против %v  (%v)", i, params, obj.Params, obj.URL)
		}
	}
}
