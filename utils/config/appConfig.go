package config

import (
	"log"
	"reflect"

	"github.com/anjolaoluwaakindipe/videoapp/utils/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AppConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
	Env  string `mapstructure:"SERVER_ENV"`
}


// Get environment variable and parse it into viper
func parseEnvIntoViper(i interface{}, viperConfig *viper.Viper, parent string, delim string) error {
	// get the type of AppConfig struct
	r := reflect.TypeOf(i)

	// Assign to a struct value if r is a pointer
	if r.Kind() == reflect.Pointer {
		r = r.Elem()
	}

	// loop through all the fields of the sturct 
	for i := 0; i < r.NumField(); i++ {
		// get the mapstructure tag from each field
		field := r.Field(i)
		env := field.Tag.Get("mapstructure")

		// if there is a parent reassign env to be parent-delimiter-child e.g DATABASE_PORT. "_" is the delimiter
		if parent != "" {
			env = parent + delim + env
		}

		// if the field contains a nested struct, recursively pass in that nested struct
		if field.Type.Kind() == reflect.Struct {
			t := reflect.New(field.Type).Elem().Interface()
			parseEnvIntoViper(t, viperConfig, env, delim)
			continue
		}

		// if no nested struct, bind the value of the mapstructure to the viper config
		if err := viperConfig.BindEnv(env); err != nil {
			// return an error if there is one
			return err
		}
	}

	// return nil if everything was successful
	return nil
}

// checks to see if appconfig has values for all its fields
func isAppConfigComplete(i interface{}, logger logger.Logger) {
	// gets the value and type of app config
	r := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	// checks if appConfig is a pointer and assigns r and v to the value of the pointer
	if r.Kind() == reflect.Pointer {
		r = r.Elem()
		v = v.Elem()
	}

	// iterate through each field of appconfig	
	for i := 0; i < r.NumField(); i++ {

		// if a field contains a nested struct call the function on that field
		if r.Field(i).Type.Kind() == reflect.Struct {
			isAppConfigComplete(v.Field(i).Interface(), logger)
		}

		// crate a zero value of that field
		emptyValue := reflect.Zero(v.Field(i).Type()).Interface()

		// check to see if any field is empty i.e the field matches its respective zero value
		if v.Field(i).Interface() == emptyValue {
			logger.Fatalf("Can't find Environment Variable: %s", r.Field(i).Tag.Get("mapstructure"))
		}
	}
}

func NewAppConfigViper(logger logger.Logger) AppConfig {
	// create an empty app config
	var ac AppConfig

	// set up the viper config instance
	viperConfig := viper.New()
	viperConfig.SetConfigFile(".env")
	viperConfig.SetConfigType("env")
	viperConfig.AddConfigPath(".")

	// if there is a .env file, load it into viperConfig
	if err := viperConfig.ReadInConfig(); err != nil {

		// if an error occured, use parseEnvIntoViper() to load environemnt variables into viper instance
		logger.Infoln("No configuration file found, using environment variables")
		viperConfig.AutomaticEnv()
		if err := parseEnvIntoViper(&ac, viperConfig, "", "_"); err != nil {
			log.Fatalln("Error parsing environment variables", err.Error())
		}
	}

	// marshal environment variables from viper config into appconfig instance
	if err := viperConfig.Unmarshal(&ac); err != nil {
		logger.Fatalln("Error while unmarshalling config/environment variables" + err.Error())
	}

	// check if all fields in appconfig have non-zero/non-default values
	isAppConfigComplete(&ac, logger)
	return ac
}

var Module = fx.Module("config", fx.Provide(NewAppConfigViper))
