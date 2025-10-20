package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer HandleError(w)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func HandleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		// log.SetOutput(os.Stderr)
		log.Println("Error:", r)
		log.Println(string(debug.Stack()))
		if err, ok := r.(Error); ok {
			WriteJSON(w, err.Code, map[string]interface{}{"message": err.String()})
		} else if err, ok := r.(error); ok {
			WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		} else {
			WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "unknown error"})
		}
	}
}

func CheckError(statusCode int, err error, customMessage ...string) {
	if err != nil {
		if len(customMessage) > 0 {
			log.Println(err)
			Abort(statusCode, errors.New(customMessage[0]))
		} else {
			Abort(statusCode, err)
		}
	}
}

type Error struct {
	Code  int
	Error interface{}
}

func (e Error) String() string {
	err := e.Error
	if bErr, ok := err.([]byte); ok {
		return string(bErr)
	} else if sErr, ok := err.(string); ok {
		return sErr
	} else if eErr, ok := err.(error); ok {
		return eErr.Error()
	} else if bErr, err := json.Marshal(err); err != nil {
		return string(bErr)
	}
	return ""
}

func Abort(code int, err interface{}) {
	if e, ok := err.(error); ok {
		r := regexp.MustCompile(`^(?:[^=]+=){2}([^-]+)`)
		newErr := r.FindAllStringSubmatch(e.Error(), -1) //will return nil when passed from a UC
		if len(newErr) < 1 || len(newErr[0]) < 2 || newErr == nil {
			panic(Error{code, errors.New(e.Error())})
		}
		panic(Error{code, errors.New(newErr[0][1])})
	} else {
		panic(Error{code, err})
	}
}
