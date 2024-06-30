package tools

import (
	"net/http"
	"os"

	"github.com/nahhoj/extensionToolsSCPI/datatypes"
	"github.com/nahhoj/extensionToolsSCPI/utils"
	"github.com/subosito/gotenv"
)

func GetSecurutyMaterial(tenant string, cookie string, credential string) datatypes.ResponseSecurutyMaterial {
	var response datatypes.ResponseSecurutyMaterial
	gotenv.Load()
	packageName := os.Getenv("packageName")
	iflowName := os.Getenv("iflowName")
	fileBase64 := os.Getenv("fileBase64")
	if tenant == "" || cookie == "" || credential == "" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "query params tenant and credential are required, header cookie is required"
		return response
	}

	scpi := utils.SCPI{
		Tenant: utils.DetectNEOFoundry(tenant),
		Cookie: cookie,
	}

	paramsCredential, err := scpi.GetSecurutyMaterial(credential)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
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

	stepTestTaskID, id, err := scpi.SimulateIflow(metadata, packageId, iflowId, credential, "securityMaterialName")
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	user, passwd, _, _, err := scpi.SimulateFinished(stepTestTaskID, id, packageId, iflowId)
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

	if paramsCredential["sec:credential.kind"] == "default" {
		response.User = user
		response.Password = passwd
	} else if paramsCredential["sec:credential.kind"] == "secure_param" {
		response.Secure = passwd
	} else {
		if paramsCredential["sec:grant.type"] == "OAuth2SAMLBearerAssertion" {
			response.ClientKey = paramsCredential["clientKey"]
		} else {
			response.ClientKey = user
			response.Secrect = passwd
		}
	}
	response.Type = paramsCredential["sec:credential.kind"]
	return response
}
