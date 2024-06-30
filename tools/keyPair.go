package tools

import (
	"net/http"
	"os"

	"github.com/nahhoj/extensionToolsSCPI/datatypes"
	"github.com/nahhoj/extensionToolsSCPI/utils"
	"github.com/subosito/gotenv"
)

func GetPairKey(tenant string, cookie string, key string) datatypes.ResponseKeyPair {
	var response datatypes.ResponseKeyPair
	gotenv.Load()
	packageName := os.Getenv("packageName")
	iflowName := os.Getenv("iflowName")
	fileBase64 := os.Getenv("fileBase64")
	if tenant == "" || cookie == "" || key == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "query params tenant and key are required, header cookie is required"
		return response
	}

	scpi := utils.SCPI{
		Tenant: utils.DetectNEOFoundry(tenant),
		Cookie: cookie,
	}

	packageId, err := scpi.CreatePackage(packageName)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	iflowId, err := scpi.UploadIflow(fileBase64, packageId, packageName, iflowName)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	metadata, err := scpi.GetMetadataIflow(packageId, iflowId)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	stepTestTaskID, id, err := scpi.SimulateIflow(metadata, packageId, iflowId, key, "keyPairName")
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	_, _, key, cert, err := scpi.SimulateFinished(stepTestTaskID, id, packageId, iflowId)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	err = scpi.DeletePackage(packageName)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}
	response.StatusCode = http.StatusOK
	response.Key = key
	response.Cert = cert
	return response
}
