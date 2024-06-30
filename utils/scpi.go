package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/nahhoj/extensionToolsSCPI/datatypes"
)

type SCPI struct {
	Tenant string
	Cookie string
}

func (scpi SCPI) GetSecurutyMaterial(credential string) (map[string]string, error) {
	var request Request
	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Url = scpi.Tenant + "/Operations/com.sap.it.km.api.commands.SecurityMaterialsListCommand"
	request.Method = "GET"
	request.Headers = headers

	response := CallHTTPService(request)
	var userCredentials datatypes.UserCredentials
	if response.StatusCode != 200 || strings.Contains(response.Body, "<html>") {
		return nil, errors.New(fmt.Sprintln("There is an error calling service GET SecurityMaterialsListCommand ", response.StatusCode, response.Body))
	}

	json.Unmarshal([]byte(response.Body), &userCredentials)
	error := fmt.Errorf("the credential %v not found", credential)
	params := make(map[string]string)
	for _, artifactInformations := range userCredentials.ArtifactInformations {
		if artifactInformations.Name == credential {
			for _, tags := range artifactInformations.Tags {
				params[tags.Name] = tags.Value
			}
			error = nil
			break
		}
	}
	return params, error
}

func (scpi SCPI) CreatePackage(packageName string) (string, error) {
	var request Request
	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Headers = headers
	request.Method = "POST"
	request.Url = scpi.Tenant + "/odata/1.0/workspace.svc/ContentEntities.ContentPackages"
	request.Body = fmt.Sprintf(`{
		"Category": "Integration",
		"SupportedPlatforms": "SAP HANA Cloud Integration",
		"TechnicalName": "%s",
		"DisplayName": "%s",
		"ShortText": "%s"
	}`, packageName, packageName, packageName)

	response := CallHTTPService(request)
	if response.StatusCode != 201 {
		if response.StatusCode == 409 {
			request = Request{}
			request.Headers = headers
			request.Method = "GET"
			request.Url = scpi.Tenant + "/odata/1.0/workspace.svc/ContentEntities.ContentPackages" + "?%24filter=TechnicalName%20eq%20%27" + packageName + "%27"
			response = CallHTTPService(request)
			if response.StatusCode != 200 {
				return "", errors.New(fmt.Sprintln("There is an error calling service GET ContentEntities.ContentPackages ", response.StatusCode, response.Body))
			}
			var contentPackage datatypes.ContentPackage
			json.Unmarshal([]byte(response.Body), &contentPackage)
			return contentPackage.D.Results[0].RegID, nil
		} else {
			return "", errors.New(fmt.Sprintln("There is an error calling service POST ContentEntities.ContentPackages ", response.StatusCode, response.Body))
		}
	} else {
		var contentPackages datatypes.ContentPackages
		json.Unmarshal([]byte(response.Body), &contentPackages)
		return contentPackages.D.RegID, nil
	}
}

func (scpi SCPI) UploadIflow(fileBase64 string, packageId string, packageName string, IflowName string) (string, error) {
	var request Request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("__xmlview11--iflowBrowse", IflowName+".zip")
	fileBase64Bytes, _ := base64.StdEncoding.DecodeString(fileBase64)
	_, err := io.WriteString(part, string(fileBase64Bytes))
	if err != nil {
		return "", errors.New(fmt.Sprintln("Cannot read file to upload"))
	}
	writer.WriteField("_charset_", "UTF-8")
	writer.WriteField("__xmlview11--iflowBrowse-data", fmt.Sprintf(`{"name":"%s","description":"","type":"IFlow","id":"%s","additionalAttrs":{"source":[],"target":[],"productProfile":["iflmap"],"nodeType":["IFLMAP"]},"packageId":"%s","fileName":"%s.zip"}`, IflowName, IflowName, packageName, IflowName))

	request.BodyReader = body
	request.Method = "POST"
	request.Url = fmt.Sprintf("%s/api/1.0/workspace/%s/iflows", scpi.Tenant, packageId)
	request.Headers = map[string][]string{
		"Cookie":       {scpi.Cookie},
		"Content-Type": {writer.FormDataContentType()},
	}
	writer.Close()
	response := CallHTTPService(request)
	if response.StatusCode != 201 {
		return "", errors.New(fmt.Sprintln("There is an error calling service POST IntegrationDesigntimeArtifacts ", response.StatusCode, response.Body))
	}
	var uploadIflow datatypes.UploadIflow
	json.Unmarshal([]byte(response.Body), &uploadIflow)
	return uploadIflow.ID, nil
}

func (scpi SCPI) GetMetadataIflow(PackageId string, iflowId string) (string, error) {
	var request Request
	request.Headers = map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Method = "GET"
	request.Url = scpi.Tenant + fmt.Sprintf("/api/1.0/workspace/%s/artifacts/%s/entities/%s/iflows/%s", PackageId, iflowId, iflowId, "Test_123")
	response := CallHTTPService(request)
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintln("There is an error calling service DELETE IntegrationPackages ", response.StatusCode, response.Body))
	}
	return response.Body, nil
}

func (scpi SCPI) SimulateIflow(metadata string, packageId string, iflowId string, credential string, property string) (string, int64, error) {
	var request Request
	id := time.Now().Unix()
	request.Headers = map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Method = "PUT"
	request.Url = scpi.Tenant + fmt.Sprintf("/api/1.0/workspace/%s/artifacts/%s/entities/%s/iflows/Test/simulations?id=%d&isReadMode=true&webdav=SIMULATE", packageId, iflowId, iflowId, id)
	request.Body = fmt.Sprintf(`{
			"startSeqID": "SequenceFlow_5",
			"endSeqID": "SequenceFlow_6",
			"process": "Process_1",
			"inputPayload": {
				"properties": {
					"%s": "%s"
				},
				"body": null
			},
			"mockPayloads": {},
			"traceCache": {},
			"iflowModelTO":%s}`, property, credential, metadata)
	response := CallHTTPService(request)
	if response.StatusCode != 200 {
		return "", id, fmt.Errorf("There is an error calling service SIMULATE simulations ", response.StatusCode, response.Body)
	}
	var simulate datatypes.Simulate
	json.Unmarshal([]byte(response.Body), &simulate)
	return simulate.StepTestTaskID, id, nil
}

func (scpi SCPI) SimulateFinished(StepTestTaskID string, id int64, packageId string, iflowId string) (string, string, string, string, error) {
	var request Request
	request.Headers = map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Method = "GET"
	request.Url = scpi.Tenant + fmt.Sprintf("/api/1.0/workspace/%s/artifacts/%s/entities/%s/iflows/%s/simulations/%s?id=%d", packageId, iflowId, iflowId, "Test_123", StepTestTaskID, id)
	var simulations datatypes.Simulations
	for {
		time.Sleep(1 * time.Second)
		response := CallHTTPService(request)
		if response.StatusCode != 200 {
			return "", "", "", "", fmt.Errorf("There is an error calling service SIMULATE simulations ", response.StatusCode, response.Body)
		}
		json.Unmarshal([]byte(response.Body), &simulations)
		if simulations.PercentageComplete == "100" {
			break
		}
	}
	return simulations.TraceData.SequenceFlow6.TracePages.Num1.Properties.User, simulations.TraceData.SequenceFlow6.TracePages.Num1.Properties.Passwd, simulations.TraceData.SequenceFlow6.TracePages.Num1.Properties.Key, simulations.TraceData.SequenceFlow6.TracePages.Num1.Properties.Cert, nil
}

func (scpi SCPI) DeletePackage(packageName string) error {
	var request Request
	request.Headers = map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {scpi.Cookie},
		"Accept":       {"application/json"},
	}
	request.Method = "DELETE"
	request.Url = fmt.Sprintf("%s/odata/1.0/workspace.svc/ContentEntities.ContentPackages('%s')", scpi.Tenant, packageName)
	response := CallHTTPService(request)
	if response.StatusCode != 204 {
		return fmt.Errorf("There is an error calling service DELETE IntegrationPackages ", response.StatusCode, response.Body)
	}
	return nil
}
