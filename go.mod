module github.com/huaweicloud/cloudeye-exporter

go 1.14

require (
	github.com/huaweicloud/golangsdk v0.0.0-20210602072215-3a6ae0cf18e8
	github.com/prometheus/client_golang v1.7.0
	github.com/prometheus/common v0.10.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
