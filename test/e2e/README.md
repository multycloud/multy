This package contains expensive integration tests that deploy infrastructure to aws and azure and check if the expected infra is deployed using command line tools.

Sample running command:

```go1.18rc1 test -v -tags e2e -run AwsKubernetes -timeout 30m -p 1 ./test/...```