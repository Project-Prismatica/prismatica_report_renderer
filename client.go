package prismatica_report_renderer

import (
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	cm     = &sync.Mutex{}
	Client PrismaticaReportRendererClient
)

func GetPrismaticaReportRendererClient() (PrismaticaReportRendererClient, error) {
	cm.Lock()
	defer cm.Unlock()

	if Client != nil {
		return Client, nil
	}

	logrus.Info("Creating prismatica_report_renderer gRPC client")
	conn, err := grpc.Dial("prismatica_report_renderer:80", grpc.DialOption(grpc.WithInsecure()))
	if err != nil {
		return nil, err
	}

	cli := NewPrismaticaReportRendererClient(conn)
	Client = cli
	return cli, nil
}
