package server

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/Project-Prismatica/prismatica_report_renderer"
	"github.com/Project-Prismatica/prismatica_report_renderer/templating_engine"

	_ "github.com/Project-Prismatica/prismatica_report_renderer/template_providers"
)

type PrismaticaReportRendererServer struct {
	prismatica_report_renderer.PrismaticaReportRendererServer
	RenderTemplateProviders []templating_engine.RenderTemplateProvider
	TemplatingEngine        *templating_engine.TemplatingEngine
}

func NewReportRenderServerOrPanic() (PrismaticaReportRendererServer) {

	/**
	 *   This is currently commented out because there is not cache eviction
	 * system in place.
	 *
	 *   Additionally, this should be moved to conform with the new template
	 * provider registration
	 *
	inMemoryTemplateCache, err := template_providers.
		NewInMemoryRenderTemplateProvider()

	if err != nil {
		logrus.Fatal("could not allocate in-memory cache")
	}
	*/

	var registeredProviders []templating_engine.RenderTemplateProvider
	for providerName, providerFactory := range
			templating_engine.RegisteredTemplateProviderFactories {
		newProvider, newProviderError := providerFactory(viper.GetViper())
		if nil != newProviderError {
			logrus.WithFields(logrus.Fields{
				"name": providerName,
				"error": newProviderError,
			}).Warn("skipping template provider")
			continue
		}
		registeredProviders = append(registeredProviders, newProvider)
	}

	createdServer := PrismaticaReportRendererServer{
		RenderTemplateProviders: registeredProviders,
		TemplatingEngine: templating_engine.NewTemplateEngine(),
	}

	return createdServer
}

func (s PrismaticaReportRendererServer) resolveTemplate(templateId string)(
		*templating_engine.ReportTemplate, bool) {
	logrus.WithFields(logrus.Fields{
		"providers": s.RenderTemplateProviders,
		"templateId": templateId,
	}).Debug("resolving template with providers")

	for _, provider := range s.RenderTemplateProviders {
		foundTemplate, err := provider.ResolveTemplate(templateId)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err,
				"provider": provider, "templateId": templateId}).
				Warn("provider error resolving template, continuing")
			continue
		}

		if foundTemplate != nil {
			return foundTemplate, true
		}
	}

	logrus.WithFields(logrus.Fields{"templateId": templateId}).
		Error("could not resolve template")
	return nil, false
}

func (s PrismaticaReportRendererServer) storeTemplate (
		template *templating_engine.ReportTemplate)(bool) {

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
