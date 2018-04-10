package template_functions

import (
	"errors"
	"fmt"

	"net/url"

	"github.com/globalsign/mgo"
	"github.com/flosch/pongo2"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

type mongoConnectionInformation struct {
	originalUri, host, database string
	authority *url.Userinfo
}

const (
	mongodbFilterName = "mongo"
)

func init() {
	pongo2.RegisterFilter(mongodbFilterName, mongodbFilterFunction)
}

func parseMongoUri(mongoUri string)(parsed mongoConnectionInformation,
		err error) {

	parsedUrl, urlError := url.Parse(mongoUri)
	if urlError != nil {
		logrus.WithField("mongoUri", mongoUri).
			Warn("could not parse mongo URL")
		err = errors.New("could not parse mongo URI")
		return
	}

	if parsedUrl.Scheme != "mongo://" {
		logrus.WithField("offendingScheme", parsedUrl.Scheme).Warn(
			"unknown scheme provided in mongo URI")
	}

	parsed.database = parsedUrl.Path
	if '/' == parsed.database[0] {
		parsed.database = parsed.database[1:]
	}

	parsed.authority = parsedUrl.User
	parsed.host = parsedUrl.Host
	parsed.originalUri = mongoUri
	return
}

func getMongoConnection(connectionInfo mongoConnectionInformation)(
		session *mgo.Session, collection *mgo.Database, err error) {
	mongoConnection, mongoConnectionError := mgo.Dial(connectionInfo.host)
	if mongoConnectionError != nil {
		logrus.WithField("error", mongoConnectionError).
			Warn("could not connect to mongo in filter")
		err = errors.New("could not connect to mongo")
		return
	}

	session = mongoConnection
	collection = session.DB(connectionInfo.database)
	return
}

func mongodbFilterFunction(mongoQueryValue *pongo2.Value,
		mongoConnectionInfoValue *pongo2.Value)(out *pongo2.Value,
		err *pongo2.Error) {
	logrus.WithFields(logrus.Fields{"mongoQuery": mongoQueryValue,
		"mongoConnectionInformation": mongoConnectionInfoValue}).Debug(
		"executing mongo query")

	filterError := pongo2.Error{Sender: mongodbFilterName}

	mongoConnectionInfo, connectionParseError := parseMongoUri(
		mongoConnectionInfoValue.String())
	if connectionParseError != nil {
		filterError.OrigError = connectionParseError
		err = &filterError
		return
	}

	mongoSession, mongoDatabase, mongoConnectionError := getMongoConnection(
		mongoConnectionInfo)
	if mongoConnectionError != nil {
		filterError.OrigError = mongoConnectionError
		err = &filterError
		return
	}
	defer mongoSession.Close()


	mongoQuery := bson.M{"eval": mongoQueryValue.String()}
	mongoOutput := bson.M{}
	mongoDatabase.Run(mongoQuery, &mongoOutput)

	if mongoOutput["ok"].(float64) != 1 {
		filterError.OrigError = errors.New(fmt.Sprintf(
			"could not execute query, '%s'", mongoOutput["errmsg"]))
		err = &filterError
		return
	}

	serializedReturnValue, jsonExtractionError := bson.MarshalJSON(
		mongoOutput["retval"])
	if jsonExtractionError != nil {
		logrus.WithField("error", jsonExtractionError).
			Warn("could not extract value")
		filterError.OrigError = errors.New("could not extract return value")
		err = &filterError
		return
	}

	out = pongo2.AsValue(string(serializedReturnValue))

	return
}
