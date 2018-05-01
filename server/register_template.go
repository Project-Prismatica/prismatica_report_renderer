package server

import (
	"errors"
	"github.com/sirupsen/logrus"

	"github.com/satori/go.uuid"
	context "golang.org/x/net/context"

	"github.com/Project-Prismatica/prismatica_report_renderer"
	"github.com/Project-Prismatica/prismatica_report_renderer/templating_engine"
)

func (s PrismaticaReportRendererServer) RegisterTemplate(ctx context.Context,
		r *prismatica_report_renderer.TemplateRegistrationRequest) (
		*prismatica_report_renderer.TemplateRegistrationResponse, error) {

	if r.RequestId == "" {
		newId, err := s.newRequestId()
		if err != nil {
			return nil, err
		}
		r.RequestId = newId
	}

	logFields := logrus.Fields{"requestId": r.RequestId}
	if logrus.GetLevel() == logrus.DebugLevel {
		logFields["renderRequest"] = r
		logFields["renderTemplate"]  = r.Template
	}

	logrus.WithFields(logFields).Info("handling template registration request")

	if r.Timestamp == nil {
		logrus.Warn("request with no timestamp provided")
		return nil, errors.New("no timestamp provided")
	}

	if len(r.Template) == 0{
		logrus.Warn("request with no template provided")
		return nil, errors.New("no template contents")
	}

	newTemplateIdSource, err := uuid.NewV4()
	if err != nil {
		logrus.Warn("could not generate new template id")
		return nil, errors.New("could not generate new template id")
	}

	newTemplateId := newTemplateIdSource.String()
	logrus.WithFields(logrus.Fields{"newTemplateId": newTemplateId}).
		Debug("generated new template Id")

	templateToStore, templateCreationError := templating_engine.
		NewTemplate(r.Template)
	if templateCreationError != nil {
		logrus.WithFields(logrus.Fields{"requestId": r.RequestId,
		"templateCreationError": templateCreationError}).
		Warn("could not load client template")
	}

	templateToStore.Id = newTemplateId
	templateToStore.Template = r.Template

	successfullyStored := s.storeTemplate(templateToStore)
	if ! successfullyStored {
		return nil, errors.New("unknown error")
	}

	response := &prismatica_report_renderer.TemplateRegistrationResponse{
		TemplateId: newTemplateId,
		RequestId: r.RequestId,
	}

	logFields = logrus.Fields{"newTemplateId": newTemplateId,
		"requestId": r.RequestId}
	if logrus.GetLevel() == logrus.DebugLevel {
		logFields["response"] = response
	}
	logrus.WithFields(logFields).Info("completed registration request")

	return response, nil
}
