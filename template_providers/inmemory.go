package template_providers

import (
	"github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"

	"github.com/Project-Prismatica/prismatica_report_renderer/templating_engine"
)

const (
	defaultCacheSize = 10
)

type InMemoryRenderTemplateProvider struct {
	storedTemplates *lru.Cache
}

func NewInMemoryRenderTemplateProvider()(*InMemoryRenderTemplateProvider, error) {
	newCache, err := lru.New(defaultCacheSize)

	if err != nil {
		return nil, err
	}

	createdProvider := &InMemoryRenderTemplateProvider{
		storedTemplates: newCache,
	}

	return createdProvider, nil
}

func (s InMemoryRenderTemplateProvider) ResolveTemplate(templateId string)(
		foundTemplate *templating_engine.ReportTemplate, err error) {

	logrus.WithFields(logrus.Fields{"templateId": templateId,
		"provider": "inMemory"}).Debug("resolving template")

	candidateTemplate, templateFound := s.storedTemplates.Get(templateId)

	if !templateFound {
		logrus.WithFields(logrus.Fields{"templateId": templateId,
			"provider": "inMemory"}).Debug("not found")

		return
	}

	foundTemplate = candidateTemplate.(*templating_engine.ReportTemplate)

	logrus.WithFields(logrus.Fields{"templateId": templateId,
		"provider": "inMemory"}).Debug("found")

	return
}

func (s *InMemoryRenderTemplateProvider) StoreTemplate(
		toStore *templating_engine.ReportTemplate)(error) {
	s.storedTemplates.Add(toStore.Id, toStore)
	return nil
}

func (s InMemoryRenderTemplateProvider) SupportsWrite()(bool) {
	return true
}
