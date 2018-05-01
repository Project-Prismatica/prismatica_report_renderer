package server

import (
	"errors"

	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"

	"github.com/Project-Prismatica/prismatica_report_renderer"
)

func (s PrismaticaReportRendererServer) Render(ctx context.Context,
		r *prismatica_report_renderer.RenderRequest) (
		*prismatica_report_renderer.RenderResponse, error) {

	logFields := logrus.Fields{"renderTemplateId": r.TemplateId}
	if logrus.GetLevel() == logrus.DebugLevel {
		logFields["renderRequest"] = r
	}

	logrus.WithFields(logFields).Info("handling render request")

	if r.Timestamp == nil {
		return nil, errors.New("no timestamp provided")
	}

	foundTemplate, templateWasFound := s.resolveTemplate(r.TemplateId)
	if ! templateWasFound {
		return nil, errors.New("unknown template Id")
	}

	logrus.WithFields(logrus.Fields{"foundTemplate": foundTemplate}).
		Debug("found template")
	out, err := s.TemplatingEngine.Render(foundTemplate, r.RenderingVariables)
	if err != nil {
		return nil, errors.New("rendering error")
	}

	response := &prismatica_report_renderer.RenderResponse{
		RequestId: r.RequestId,
		Result: out,
	}
	return response, nil
}
