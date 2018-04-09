package server

import (
	"errors"
	"github.com/satori/go.uuid"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/Project-Prismatica/prismatica_report_renderer"
)

type PrismaticaReportRendererServer struct {
	prismatica_report_renderer.PrismaticaReportRendererServer
	RenderTemplateProviders []prismatica_report_renderer.RenderTemplateProvider
	TemplatingEngine prismatica_report_renderer.TemplatingEngine
}

func NewReportRenderServerOrPanic() (PrismaticaReportRendererServer) {

	inMemoryTemplateCache, err := prismatica_report_renderer.
		NewInMemoryRenderTemplateProvider()

	if err != nil {
		logrus.Fatal("could not allocate in-memory cache")
	}

	createdServer := PrismaticaReportRendererServer{
		RenderTemplateProviders: []prismatica_report_renderer.
			RenderTemplateProvider {
				inMemoryTemplateCache,
		},
	}

	return createdServer
}

func (s PrismaticaReportRendererServer) resolveTemplate(templateId string)(
		*prismatica_report_renderer.ReportTemplate, bool) {
	logrus.WithFields(logrus.Fields{"templateId": templateId}).
		Debug("resolving template with providers")

	for _, provider := range s.RenderTemplateProviders {
		foundTemplate, err := provider.ResolveTemplate(templateId)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err,
				"provider": provider, "templateId": templateId}).
				Error("provider error resolving template, continuing")
			continue
		}

		if foundTemplate != nil {
			return foundTemplate, true
		}
	}

	logrus.WithFields(logrus.Fields{"templateId": templateId}).
		Info("could not resolve template")
	return nil, false
}

func (s PrismaticaReportRendererServer) storeTemplate (
		template *prismatica_report_renderer.ReportTemplate)(bool) {

	everStoredSuccessfully := false

	for _, provider := range s.RenderTemplateProviders {
		if ! provider.SupportsWrite() {
			logrus.WithFields(logrus.Fields{
				"providerType": reflect.TypeOf(provider).String()}).
				Debug("provider does not support write, skipping")
			continue
		}

		err := provider.StoreTemplate(template)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"providerType": reflect.TypeOf(provider).String(),
				"error": err}).
				Error("provider could not save template")
			continue
		}

		everStoredSuccessfully = true
	}

	return everStoredSuccessfully
}

func (s PrismaticaReportRendererServer) newRequestId() (requestId string,
		err error) {

	newRequestIdUuid, uuidCreationError := uuid.NewV4()
	if uuidCreationError != nil {
		logrus.WithField("error", uuidCreationError).
			Error("could not create new request Id")
		err = errors.New("could not allocate request id")
		return
	}
	requestId = newRequestIdUuid.String()

	return
}
