package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"

{{renderImports .ServerImports}}

	"github.com/mimokpl/kratos-bootstrap/bootstrap"
	conf "github.com/mimokpl/kratos-bootstrap/api/gen/go/conf/v1"

	//_ "github.com/mimokpl/kratos-bootstrap/config/apollo"
	//_ "github.com/mimokpl/kratos-bootstrap/config/consul"
	//_ "github.com/mimokpl/kratos-bootstrap/config/etcd"
	//_ "github.com/mimokpl/kratos-bootstrap/config/kubernetes"
	//_ "github.com/mimokpl/kratos-bootstrap/config/nacos"
	//_ "github.com/mimokpl/kratos-bootstrap/config/polaris"

	//_ "github.com/mimokpl/kratos-bootstrap/logger/aliyun"
	//_ "github.com/mimokpl/kratos-bootstrap/logger/fluent"
	//_ "github.com/mimokpl/kratos-bootstrap/logger/logrus"
	//_ "github.com/mimokpl/kratos-bootstrap/logger/tencent"
	//_ "github.com/mimokpl/kratos-bootstrap/logger/zap"
	//_ "github.com/mimokpl/kratos-bootstrap/logger/zerolog"

	//_ "github.com/mimokpl/kratos-bootstrap/registry/consul"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/etcd"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/eureka"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/kubernetes"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/nacos"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/polaris"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/servicecomb"
	//_ "github.com/mimokpl/kratos-bootstrap/registry/zookeeper"

	//_ "github.com/mimokpl/kratos-bootstrap/tracer"

	"{{.Module}}/pkg/serviceid"
)

var version = "1.0.0"

// go build -ldflags "-X main.version=x.y.z"

func newApp(
	ctx *bootstrap.Context,
{{renderFormalParameters .ServerFormalParameters}}) *kratos.App {
	return bootstrap.NewApp(ctx,
{{renderInParameters .ServerTransferParameters 2}}	)
}

func runApp() error {
	ctx := bootstrap.NewContext(
		context.Background(),
		&conf.AppInfo{
			Project: serviceid.ProjectName,
			AppId:   serviceid.{{renderServiceName .Service}},
			Version: version,
		},
	)
	return bootstrap.RunApp(ctx, initApp)
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}
