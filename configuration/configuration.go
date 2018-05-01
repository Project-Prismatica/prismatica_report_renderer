package configuration

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	serviceEnvPrefix = "PRISMATICA_"
	ConfigKeyLogLevelVerbose = "VERBOSE"
	ConfigKeyLogLevelDebug = "DEBUG"
	ConfigKeyMongodbUri = "MONGODB_URI"
)

type ProgramConfiguration struct {
	LogLevel	logrus.Level
	MongodbURI	string
}

func init()  {
	viper.SetEnvPrefix(serviceEnvPrefix)
	viper.AutomaticEnv()

	viper.SetDefault(ConfigKeyLogLevelVerbose, false)
	viper.SetDefault(ConfigKeyLogLevelDebug, false)
}
