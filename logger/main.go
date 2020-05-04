package logger

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/apex/log/handlers/es"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
	"github.com/tj/go-elastic"
	//"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"os"
)

var (
	level string
)

func InitLogger() {

	err := viper.BindEnv("log.level")

	if err != nil {
		fmt.Println("Bind env error %s", err)
	}

	level = viper.GetString("log.level")
	log.SetLevelFromString(level)

	handlers := viper.GetStringSlice("log.handlers")

	if len(handlers) > 1 {
		handlersArr := make([]log.Handler, len(handlers))
		for i, name := range handlers {
			h, err := getHandler(name)

			if err != nil {
				fmt.Errorf("Can`t init log handler:%s \nError:%s\n", name, err)
				continue
			}
			handlersArr[i] = h
		}

		if len(handlersArr) < 1 {
			h, err := getHandler("cli")
			if err != nil {
				fmt.Errorf("Can`t init log handler:%s \nError:%s\n", "cli", err)
			}

			log.SetHandler(h)

		} else {
			log.SetHandler(multi.New(handlersArr...))
		}

	} else {
		h, err := getHandler(handlers[0])

		if err != nil {
			fmt.Errorf("Can`t init log handler:%s \nError:%s\n", handlers[0], err)
		}

		log.SetHandler(h)

	}
}

func getHandler(name string) (log.Handler, error) {

	switch name {

	case "file":

		if !viper.IsSet("log.config.file.path") {
			return nil, errors.New("No file path set [logs.config.file.path]")
		}

		fileName := viper.GetString("log.config.file.path")
		fileFormat := "text"

		if viper.IsSet("log.config.file.format") {
			fileFormat = viper.GetString("log.config.file.format")
		}

		switch fileFormat {
		case "json":

			file, err := getFile(fileName)
			if err != nil {
				return nil, err
			}

			return json.New(file), nil

		case "text":
			file, err := getFile(fileName)

			if err != nil {
				return nil, err
			}

			return text.New(file), nil
		}

	case "cli":
		return cli.New(os.Stdout), nil
	case "es":

		if viper.IsSet("log.config.es") {

			esDSN := fmt.Sprintf("%s://%s:%d",
				viper.GetString("log.config.es.schema"),
				viper.GetString("log.config.es.host"),
				viper.GetInt("log.config.es.port"),
			)

			esTimeout, _ := time.ParseDuration("5s")
			if viper.IsSet("log.config.es.timeout") {
				esTimeout = viper.GetDuration("log.config.es.timeout")
			}

			esBuffer := 100
			if viper.IsSet("log.config.es.buffer") {
				esBuffer = viper.GetInt("log.config.es.buffer")
			}

			esClient := elastic.New(esDSN)
			esClient.HTTPClient = &http.Client{
				Timeout: esTimeout,
			}

			return es.New(&es.Config{
				Client:     esClient,
				BufferSize: esBuffer,
			}), nil
		}
		return nil, errors.New("No elastic config [log.config.es]")
	}

	return text.New(os.Stdout), nil
}

func getFile(fileName string) (io.Writer, error) {

	workDir, errDir := os.Getwd()
	if errDir != nil {
		return nil, errDir
	}

	filePath := path.Join(workDir, fileName)

	var _, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		var file, errCreate = os.Create(filePath)
		if errCreate != nil {
			fmt.Errorf("Can`t create log file:%s \nError:%s\n", filePath, errCreate)
			return nil, err
		}
		return file, nil
	}

	var file, errOpen = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errOpen != nil {
		fmt.Errorf("Can`t open log file:%s \nError:%s\n", filePath, errOpen)
		return nil, err
	}

	return file, nil
}
