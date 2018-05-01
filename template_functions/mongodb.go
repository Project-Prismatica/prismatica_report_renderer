package template_functions

import (
	"errors"
	"fmt"

	"github.com/flosch/pongo2"
	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2/bson"

	"github.com/Project-Prismatica/prismatica_report_renderer/util"
)

const (
	mongodbFilterName = "mongo"
)

func init() {
	pongo2.RegisterFilter(mongodbFilterName, mongodbFilterFunction)
}

func mongodbFilterFunction(mongoQueryValue *pongo2.Value,
		mongoConnectionInfoValue *pongo2.Value)(out *pongo2.Value,
		err *pongo2.Error) {
	logrus.WithFields(logrus.Fields{"mongoQuery": mongoQueryValue,
		"MongodbConnectionInformation": mongoConnectionInfoValue}).Debug(
		"executing mongo query")

	filterError := pongo2.Error{Sender: mongodbFilterName}

	mongoConnectionInfo, connectionParseError := util.ParseMongodbUri(
		mongoConnectionInfoValue.String())
	if connectionParseError != nil {
		filterError.OrigError = connectionParseError
		err = &filterError
		return
	}

	mongoSession, mongoDatabase, mongoConnectionError := util.GetMongodbConnection(
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
