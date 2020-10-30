package uexec

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetReportCaller(false)
}

// ErrorHandler struct for handling errors
type ErrorHandler struct {
	Logger       logrus.FieldLogger
	action       Action
	callBackFunc interface{}
	callBackArgs []interface{}
	enableGCB    bool
	elasticity   bool
	onErAction   string
	erP          int
}

// NewErrorHandler factory method to create a new errorhandler
func NewErrorHandler() *ErrorHandler {
	errorHandler := &ErrorHandler{
		Logger: log,
		action: Action{},
		erP:    -1,
	}

	return errorHandler
}

// AddGenericCallBack add a generic callback function that can run on Exec errors
func (e *ErrorHandler) AddGenericCallBack(callBackFunc interface{}, callBackArgs ...interface{}) *ErrorHandler {
	e.callBackFunc = callBackFunc
	e.callBackArgs = callBackArgs
	e.enableGenericCallBack(true)
	log.WithFields(logrus.Fields{
		"Action": "AddGenericCallBack",
		"To":     "-",
	}).Debug("Generic CallBack added successfully")

	return e
}

// DelGenericCallBack deletes the generic callback function that can run on Exec errors
func (e *ErrorHandler) DelGenericCallBack() *ErrorHandler {
	e.callBackFunc = nil
	e.callBackArgs = nil
	e.enableGenericCallBack(false)
	log.WithFields(logrus.Fields{
		"Action": "DelGenericCallBack",
		"To":     "-",
	}).Debug("Generic CallBack deleted successfully")

	return e
}

// CallBack function for running after an error has been caught
func (e *ErrorHandler) CallBack() *ErrorHandler {
	switch fn := e.callBackFunc.(type) {
	case func(...interface{}):
		fn(e.callBackArgs)
	}
	return e
}

// enableGenericCallBack for enabling / disabling globacl callback function
func (e *ErrorHandler) enableGenericCallBack(l bool) *ErrorHandler {
	e.enableGCB = l
	log.WithFields(logrus.Fields{
		"Action": "EnableReportCaller",
		"To":     l,
	}).Debug("ReportCaller has been changed")
	return e
}

// SetElasticity for changing the Elasticity of ErrorHandler
func (e *ErrorHandler) SetElasticity(el bool) *ErrorHandler {
	e.elasticity = el
	log.WithFields(logrus.Fields{
		"Action": "SetElasticity",
		"To":     el,
	}).Debug("SetElasticity has been changed")
	return e
}

// SetLogLevel for chaning the loglevel of logrus
func (e *ErrorHandler) SetLogLevel(l string) *ErrorHandler {
	switch l {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
	log.WithFields(logrus.Fields{
		"Action": "SetLogLevel",
		"To":     l,
	}).Debug("LogLevel has been changed")
	return e
}

// EnableReportCaller for enabling / disabling reportcaller of logrus
func (e *ErrorHandler) EnableReportCaller(l bool) *ErrorHandler {
	log.SetReportCaller(l)
	log.WithFields(logrus.Fields{
		"Action": "EnableReportCaller",
		"To":     l,
	}).Debug("ReportCaller has been changed")
	return e
}

// ErP set explicit index of error in the Values []interface{}
func (e *ErrorHandler) ErP(i int) *ErrorHandler {
	e.erP = i
	return e
}

// OnErr sets the action to be taken if Exec returns with error
func (e *ErrorHandler) OnErr(a string) *ErrorHandler {
	switch a {
	case "exit":
		e.onErAction = a
	case "callback":
		if !e.enableGCB {
			e.onErAction = "exit"
			log.WithFields(logrus.Fields{
				"Action": "OnErr",
				"To":     "exit",
			}).Error("Has not been changed. Please add a callback function before you do this")
			log.WithFields(logrus.Fields{
				"Action": "OnErr",
				"To":     "exit",
			}).Error("OnErr is set to exit. Any error will cause a program exit 1")
			return e
		}
		e.onErAction = a
	default:
		e.onErAction = "exit"
	}
	log.WithFields(logrus.Fields{
		"Action": "OnErr",
		"To":     a,
	}).Debug("OnErr has been changed")
	return e
}

// Exec factory for creating a new action
func (e *ErrorHandler) Exec(cmd ...interface{}) Action {
	e.action = Action{}

	if len(cmd) == 0 {
		return e.action
	}

	if len(cmd)-1 < e.erP {
		log.WithFields(logrus.Fields{
			"Action": "Exec",
			"Status": "-",
		}).Error("The Error Pointer, points to a none existant cmd index")
		os.Exit(1)
	}

	e.action.CallBackArgs = append(e.action.CallBackArgs, cmd...)
	e.getErr(cmd...)
	e.getValues(cmd...)
	e.checkE()

	e.ErP(-1)
	return e.action
}

func (e *ErrorHandler) getErr(cmd ...interface{}) {
	if e.erP >= 0 {
		e.action.Err = cmd[e.erP]
		return
	}

	for i, j := range cmd {
		if isErr, err := e.ifErr(j); isErr {
			e.action.Err = err
			e.erP = i
			return
		}
	}
	e.action.Err = nil
}

// ifErr returns true/false and the error *if any* from an interface
func (e *ErrorHandler) ifErr(x interface{}) (bool, error) {
	switch x.(type) {
	case error:
		return true, x.(error)
	default:
		return false, nil
	}
}

func (e *ErrorHandler) getValues(cmd ...interface{}) {
	for i, j := range cmd {
		if i == e.erP {
			continue
		}
		e.action.Values = append(e.action.Values, j)
	}
}

// CheckE for checking the an error result and exiting if it is not nil
func (e *ErrorHandler) checkE() {
	if e.erP < 0 || e.action.Err == nil {
		return
	}

	e.Logger.Error(e.action.Err)
	e.errAction()
}

func (e *ErrorHandler) errAction() {
	if e.onErAction == "exit" {
		if !e.elasticity {
			os.Exit(1)
		}
	}

	if e.onErAction == "callback" {
		e.CallBack()
	}
}
