package browse

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/portainer/agent"
	"github.com/portainer/agent/filesystem"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type browseRenamePayload struct {
	CurrentFilePath string
	NewFilePath     string
}

func (payload *browseRenamePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.CurrentFilePath) {
		return agent.Error("Current file path is invalid")
	}
	if govalidator.IsNull(payload.NewFilePath) {
		return agent.Error("New file path is invalid")
	}
	return nil
}

// PUT request on /browse/rename?id=:id
func (handler *Handler) browseRename(rw http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload browseRenamePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	volumeID, err := request.RetrieveQueryParameter(r, "volumeID", true)
	if volumeID != "" {
		payload.CurrentFilePath, err = filesystem.BuildPathToFileInsideVolume(volumeID, payload.CurrentFilePath)
		if err != nil {
			return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume", err}
		}
		payload.NewFilePath, err = filesystem.BuildPathToFileInsideVolume(volumeID, payload.NewFilePath)
		if err != nil {
			return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume", err}
		}
	}

	err = filesystem.RenameFile(payload.CurrentFilePath, payload.NewFilePath)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to rename file", err}
	}

	return response.Empty(rw)
}

// PUT request on /v1/browse/rename?id=:id
func (handler *Handler) browseRenameV1(rw http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	volumeID, err := request.RetrieveRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume identifier route variable", err}
	}

	var payload browseRenamePayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	payload.CurrentFilePath, err = filesystem.BuildPathToFileInsideVolume(volumeID, payload.CurrentFilePath)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume", err}
	}
	payload.NewFilePath, err = filesystem.BuildPathToFileInsideVolume(volumeID, payload.NewFilePath)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid volume", err}
	}

	err = filesystem.RenameFile(payload.CurrentFilePath, payload.NewFilePath)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to rename file", err}
	}

	return response.Empty(rw)
}