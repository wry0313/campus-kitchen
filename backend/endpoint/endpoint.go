package endpoint

import (
	server_error "backend/error"
	"backend/validator"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


func WriteWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Failed to encode API response into JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

type errResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WriteWithError sets the response header to application/json, writes the header
// with a status code, and returns an error response with a status and mesage
func WriteWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := errResponse{
		Status:  statusCode,
		Message: errMsg,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Printf("Failed to encode error response into JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// buildDecodeErrorMsg formats the decode error and returns it as a string
func buildDecodeErrorMsg(field string, want string, got string) string {
	return fmt.Sprintf("Expected %s to be %s, got %s", field, want, got)
}

// HandleDecodeErr responds with the appropriate decode error msg and sets
// the http status to 400
func HandleDecodeErr(w http.ResponseWriter, err error) {
	errMsg := server_error.ErrMsgJSONDecode
	if err, ok := err.(*json.UnmarshalTypeError); ok {
		errMsg = buildDecodeErrorMsg(err.Field, err.Type.String(), err.Value)
	}

	WriteWithError(w, http.StatusBadRequest, errMsg)
}

// WriteValidationErr responds with the appropriate validation error msg and
// sets the http status to 400
func WriteValidationErr(w http.ResponseWriter, s interface{}, err error) {
	errMsg := server_error.ErrMsgInvalidReq
	validationErrMsg := validator.GetValidationErrMsg(s, err)
	if validationErrMsg != "" {
		errMsg = validationErrMsg
	}

	WriteWithError(w, http.StatusBadRequest, errMsg)
}

func DecodeAndValidateJson(w http.ResponseWriter, r *http.Request, validator *validator.Validate, output interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(output); err != nil {
		log.Printf("api: failed to decode request: %v", err)
		HandleDecodeErr(w, err)
		return err
	}

	if err := validator.Struct(output); err != nil {
		log.Printf("api: failed to validate input: %v", err)
		WriteValidationErr(w, output, err)
		return err
	}
	return nil
}
