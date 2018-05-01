package util

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
)

type MongodbConnectionInformation struct {
	OriginalUri, Host, Database, Collection string
	UriPathParts                            []string
	Authority                               *url.Userinfo
}

func ParseMongodbUri(mongoUri string)(parsed MongodbConnectionInformation,
	err error) {

	parsedUrl, urlError := url.Parse(mongoUri)
	if urlError != nil {
		logrus.WithField("mongoUri", mongoUri).
			Warn("could not parse mongo URL")
		err = errors.New("could not parse mongo URI")
		return
	}

	if "mongo" != parsedUrl.Scheme && "mongodb" != parsedUrl.Scheme{
		logrus.WithField("offendingScheme", parsedUrl.Scheme).Warn(
			"unknown scheme provided in mongo URI")
	}

	cleanPath := path.Clean(parsedUrl.Path)
	pathParts := strings.Split(cleanPath, "/")
	if 0 == len(pathParts) {
		// TODO
	}

	for 0 == len(pathParts[0]) {
		pathParts = pathParts[1:]
	}

	parsed.Database = pathParts[0]

	if 2 == len(pathParts) {
		parsed.Collection = pathParts[1]
	}
	parsed.UriPathParts = pathParts

	parsed.Authority = parsedUrl.User
	parsed.Host = parsedUrl.Host
	parsed.OriginalUri = mongoUri
	return
}

func GetMongodbConnection(connectionInfo MongodbConnectionInformation)(
	session *mgo.Session, collection *mgo.Database, err error) {
	mongoConnection, mongoConnectionError := mgo.Dial(connectionInfo.Host)
	if mongoConnectionError != nil {
		logrus.WithField("error", mongoConnectionError).
			Warn("could not connect to mongo in filter")
		err = errors.New("could not connect to mongo")
		return
	}

	session = mongoConnection
	collection = session.DB(connectionInfo.Database)
	return
}

