package template_providers

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Project-Prismatica/prismatica_report_renderer/templating_engine"
	"github.com/Project-Prismatica/prismatica_report_renderer/util"
)

type MongodbTemplateProvider struct {
	collection			*mgo.Collection
	connectionInfo		util.MongodbConnectionInformation
	database			*mgo.Database
	session				*mgo.Session
}

type TemplateStoredInMongodb struct {
	Id	 string `bson:"_id"`
	Name string `bson:"name"`
	Text string `nson:"text"`
}

const (
	ProviderName = "MongodbTemplateProvider"
	ConfigurationKeyDefaultTemplateCollectionUri = "TEMPLATE_COLLECTION_URI"
	ConfigurationValueDefaultTemplateCollectionUri =
		"mongodb://localhost:3001/meteor/reports.ReportTemplates"
)

func init() {
	templating_engine.RegisterTemplateProvider(ProviderName,
		MongodbTemplateProviderFactory)
	viper.SetDefault(ConfigurationKeyDefaultTemplateCollectionUri,
		ConfigurationValueDefaultTemplateCollectionUri)
}

func (s TemplateStoredInMongodb) String () string {
	cutDownText := s.Text[:]
	if 100 < len(cutDownText) {
		cutDownText = cutDownText[:100]
	}
	return fmt.Sprintf(
		"TemplateStoredInMongodb{id=%s name='%s' text='%s'}",
		s.Id, s.Name, cutDownText,
	)
}

func MongodbTemplateProviderFactory(configurationSource *viper.Viper) (
		createdProvider templating_engine.RenderTemplateProvider, err error) {
	mongodbUri := configurationSource.GetString(
		ConfigurationKeyDefaultTemplateCollectionUri)

	connectionInfo, connectionParseError := util.ParseMongodbUri(mongodbUri)
	if nil != err {
		err = connectionParseError
		return
	}

	logrus.WithFields(logrus.Fields{
		"mongodbHost": connectionInfo.Host,
		"database": connectionInfo.Database,
		"collection": connectionInfo.Collection,
	}).Info("connecting to mongodb")


	session, database, mongodbConnectionError := util.GetMongodbConnection(
		connectionInfo)
	if nil != mongodbConnectionError {
		err = mongodbConnectionError
		return
	}

	logrus.WithField("mongodbHosts", session.LiveServers()).
		Info("connected to mongodb")

	result := new(MongodbTemplateProvider)
	result.connectionInfo = connectionInfo
	result.session = session
	result.database = database
	result.collection = database.C(connectionInfo.Collection)

	createdProvider = result
	return
}

func (s *MongodbTemplateProvider) ResolveTemplate(templateId string) (
		foundTemplate *templating_engine.ReportTemplate, err error) {

	templateQuery := bson.M{"_id": templateId}
	fetchedTemplate := TemplateStoredInMongodb{}
	queryError := s.collection.Find(templateQuery).One(&fetchedTemplate)
	if nil != queryError {
		logrus.WithFields(logrus.Fields{
			"templateId": templateId,
			"mongodbQueryError": queryError,
		}).Warn("could not fetch template")
		err = errors.New(fmt.Sprintf(
			"could not fetch template, %s", queryError))
		return
	}

	logrus.WithFields(logrus.Fields{"templateId": templateId,
		"fetchedTemplate": fetchedTemplate}).Debug("tried to resolve")

	instantiatedTemplate, templateInstantiationError := templating_engine.
		NewTemplate(fetchedTemplate.Text)
	if nil != templateInstantiationError {
		logrus.WithFields(logrus.Fields{
			"instantiationError":  templateInstantiationError,
			"offendingTemplateId": templateId,
		}).Warn("could not instantiate template")
		err = errors.New("could not instantiate template")
		return
	}
	instantiatedTemplate.Id = templateId

	foundTemplate = instantiatedTemplate
	return
}

func (s *MongodbTemplateProvider) Shutdown() () {
	s.session.Close()
}

func (s *MongodbTemplateProvider) StoreTemplate(
		toStore *templating_engine.ReportTemplate) (err error) {
	panic("attempted to write to the mongodb template provider, it doesn't" +
		" support writing")
	return
}

func (s MongodbTemplateProvider) SupportsWrite() (bool) {
	return false
}

func (s MongodbTemplateProvider) String() string {
	return fmt.Sprintf("MongodbTemplateProvider{uri:%s}",
		s.connectionInfo.OriginalUri)
}
