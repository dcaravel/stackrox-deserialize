module github.com/dcaravel/stackrox-deserialize

go 1.23.4

toolchain go1.23.7

replace github.com/stackrox/rox => github.com/stackrox/stackrox v0.0.0-20250328184422-41a2f7649e51

require github.com/stackrox/rox v0.0.0-00010101000000-000000000000

require (
	github.com/planetscale/vtprotobuf v0.6.1-0.20240409071808-615f978279ca // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
