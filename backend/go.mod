module github.com/vasiliiperfilev/cookie

go 1.20

require (
	github.com/gorilla/websocket v1.5.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.9.0
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1
)

replace (
	github.com/cucumber/godog => github.com/laurazard/godog v0.0.0-20220922095256-4c4b17abdae7

	golang.org/x/oauth2 => golang.org/x/oauth2 v0.1.0

	// For k8s dependencies, we use a replace directive, to prevent them being
	// upgraded to the version specified in containerd, which is not relevant to the
	// version needed.
	// See https://github.com/docker/buildx/pull/948 for details.
	// https://github.com/docker/buildx/blob/v0.8.1/go.mod#L62-L64
	k8s.io/api => k8s.io/api v0.22.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.4
	k8s.io/client-go => k8s.io/client-go v0.22.4
)
