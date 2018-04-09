package server

import (
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/Project-Prismatica/prismatica_report_renderer"
	"github.com/lileio/lile"
)

var (
	s = NewReportRenderServerOrPanic()
	cli prismatica_report_renderer.PrismaticaReportRendererClient
)

func TestMain(m *testing.M) {
	impl := func(g *grpc.Server) {
		prismatica_report_renderer.RegisterPrismaticaReportRendererServer(g, s)
	}

	gs := grpc.NewServer()
	impl(gs)

	addr, serve := lile.NewTestServer(gs)
	go serve()

	cli = prismatica_report_renderer.NewPrismaticaReportRendererClient(lile.TestConn(addr))

	os.Exit(m.Run())
}
