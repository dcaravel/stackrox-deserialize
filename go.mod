module github.com/dcaravel/stackrox-deserialize

go 1.25.0

replace github.com/stackrox/rox => github.com/stackrox/stackrox v0.0.0-20260204210318-9aeddf591777

require (
	github.com/stackrox/rox v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240409071808-615f978279ca // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
