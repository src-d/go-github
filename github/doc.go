// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package github provides a client for using the GitHub API.

Usage:

	import "github.com/src-d/go-github/github"

Construct a new GitHub client, then use the various services on the client to
access different parts of the GitHub API. For example:

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List("willnorris", nil)

Some API methods have optional parameters that can be passed. For example:

	client := github.NewClient(nil)

	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg("github", opt)

The services of a client divide the API into logical chunks and correspond to
the structure of the GitHub API documentation at
http://developer.github.com/v3/.

Authentication

The go-github library does not directly handle authentication. Instead, when
creating a new client, pass an http.Client that can handle authentication for
you. The easiest and recommended way to do this is using the golang.org/x/oauth2
library, but you can always use any other library that provides an http.Client.
If you have an OAuth2 access token (for example, a personal API token), you can
use it with the oauth2 library using:

	import "golang.org/x/oauth2"

	func main() {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "... your access token ..."},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)

		// list all repositories for the authenticated user
		repos, _, err := client.Repositories.List("", nil)
	}

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.

See the oauth2 docs for complete instructions on using that library.

For API methods that require HTTP Basic Authentication, use the
BasicAuthTransport.

Rate Limiting

GitHub imposes a rate limit on all API clients. Unauthenticated clients are
limited to 60 requests per hour, while authenticated clients can make up to
5,000 requests per hour. To receive the higher rate limit when making calls
that are not issued on behalf of a user, use the
UnauthenticatedRateLimitedTransport.

The Rate method on a client returns the rate limit information based on the most
recent API call. This is updated on every call, but may be out of date if it's
been some time since the last API call and other clients have made subsequent
requests since then. You can always call RateLimits() directly to get the most
up-to-date rate limit data for the client.

To detect an API rate limit error, you can check if its type is *github.RateLimitError:

	repos, _, err := client.Repositories.List("", nil)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
	}

Learn more about GitHub rate limiting at
http://developer.github.com/v3/#rate-limiting.

Conditional Requests

The GitHub API has good support for conditional requests which will help
prevent you from burning through your rate limit, as well as help speed up your
application. go-github does not handle conditional requests directly, but is
instead designed to work with a caching http.Transport. We recommend using
https://github.com/gregjones/httpcache for that.

Learn more about GitHub conditional requests at
https://developer.github.com/v3/#conditional-requests.

Creating and Updating Resources

All structs for GitHub resources use pointer values for all non-repeated fields.
This allows distinguishing between unset fields and those set to a zero-value.
Helper functions have been provided to easily create these pointers for string,
bool, and int values. For example:

	// create a new private repository named "foo"
	repo := &github.Repository{
		Name:    github.String("foo"),
		Private: github.Bool(true),
	}
	client.Repositories.Create("", repo)

Users who have worked with protocol buffers should find this pattern familiar.

Pagination

All requests for resource collections (repos, pull requests, issues, etc.)
support pagination. Pagination options are described in the
github.ListOptions struct and passed to the list methods directly or as an
embedded type of a more specific list options struct (for example
github.PullRequestListOptions). Pages information is available via the
github.Response struct.

	client := github.NewClient(nil)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg("github", opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

*/
package github
