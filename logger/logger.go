package logger

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"time"
)

type Params map[string]interface{}

func InitLogger(dsn string) error {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	if err := InitSentry(dsn); err != nil {
		return err
	}
	return nil
}

type Msg struct {
	Message string
	Params  map[string]interface{}
}

func (msg *Msg) String() string {
	if len(msg.Params) > 0 {
		return fmt.Sprintf("%s - %v", msg.Message, msg.Params)
	}
	return fmt.Sprintf(msg.Message)
}

func Info(args ...interface{}) {
	if len(args) == 0 {
		Panic("call to logger.Info with no arguments")
	}
	msg := getMessage(args...)
	log.WithFields(msg.Params).Info(msg.Message)
}

func Debug(args ...interface{}) {
	if len(args) == 0 {
		Panic("call to logger.Debug with no arguments")
	}
	msg := getMessage(args...)
	log.WithFields(msg.Params).Debug(msg.Message)
}

func Warn(args ...interface{}) {
	if len(args) == 0 {
		Panic("call to logger.Warn with no arguments")
	}
	msg := getMessage(args...)
	log.WithFields(msg.Params).Warn(msg.Message)
}

func getMessage(args ...interface{}) *Msg {
	msg := &Msg{Params: make(Params)}
	var generic []string
	var message []string
	for _, arg := range args {
		switch arg := arg.(type) {
		case nil:
			continue
		case string:
			message = append(message, arg)
		case Params:
			appendMap(msg.Params, arg)
		case map[string]interface{}:
			appendMap(msg.Params, arg)
		default:
			generic = append(generic, fmt.Sprintf("%v", arg))
		}
	}
	if len(message) > 0 {
		msg.Message = strings.Join(message[:], ": ")
	}
	if len(generic) > 0 {
		msg.Params["objects"] = strings.Join(generic[:], " | ")
	}
	return msg
}

func appendMap(root map[string]interface{}, tmp map[string]interface{}) {
	for k, v := range tmp {
		root[k] = v
	}
}

func LogRequest(stop time.Duration, currency, request string) {
	if stop > (time.Second * 2) {
		Error("Response time exception", Params{
			"currency": currency,
			"request":  request,
			"time":     stop.String(),
		})
	} else {
		SendMessage(request + ": currency - " + currency + " time - " + stop.String())
	}
}
