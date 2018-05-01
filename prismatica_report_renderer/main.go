package main

import (
	"github.com/Project-Prismatica/prismatica_report_renderer"
	"github.com/Project-Prismatica/prismatica_report_renderer/prismatica_report_renderer/cmd"
	"github.com/Project-Prismatica/prismatica_report_renderer/server"
	"github.com/lileio/lile"
	"github.com/lileio/lile/fromenv"
	"github.com/lileio/pubsub"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	s := server.NewReportRenderServerOrPanic()

	lile.Name("prismatica_report_renderer")
	lile.Server(func(g *grpc.Server) {
		prismatica_report_renderer.RegisterPrismaticaReportRendererServer(g, s)
	})

	pubsub.SetClient(&pubsub.Client{
		ServiceName: lile.GlobalService().Name,
		Provider:    fromenv.PubSubProvider(),
	})

	cmd.Execute(viper.GetViper())
}
