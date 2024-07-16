package datatypes

//Data types service provider

type ResponseSecurutyMaterial struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Message    string `json:"message,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	ClientKey  string `json:"clientKey,omitempty"`
	Secrect    string `json:"secrect,omitempty"`
	Secure     string `json:"secure,omitempty"`
	Type       string `json:"type,omitempty"`
}

type ResponseKeyPair struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Message    string `json:"message,omitempty"`
	Key        string `json:"key,omitempty"`
	Cert       string `json:"cert,omitempty"`
}

type GroovyLog struct {
	Script     string `json:"script,omitempty"`
	Body       string `json:"body,omitempty"`
	Headers    string `json:"headers,omitempty"`
	Properties string `json:"properties,omitempty"`
	Method     string `json:"method,omitempty"`
	Log        string `json:"log,omitempty"`
}

type FormatCode struct {
	Code string `json:"code"`
}

//Data types service consumer API SAP CPI
type UserCredentials struct {
	ArtifactInformations []struct {
		Name string `json:"name"`
		Tags []struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"tags"`
	} `json:"artifactInformations"`
}

type ContentPackages struct {
	D struct {
		RegID     string `json:"reg_id"`
		Artifacts struct {
			Results []struct {
				Name  string `json:"Name"`
				Type  string `json:"Type"`
				RegID string `json:"reg_id"`
			} `json:"results"`
		} `json:"Artifacts"`
	} `json:"d"`
}

type ContentPackage struct {
	D struct {
		Results []struct {
			RegID string `json:"reg_id"`
		} `json:"results"`
	} `json:"d"`
}

type UploadIflow struct {
	ID string `json:"id"`
}

type Simulate struct {
	StepTestTaskID string `json:"stepTestTaskId"`
}

type Simulations struct {
	StatusCode         string `json:"statusCode"`
	StatusMessage      string `json:"statusMessage"`
	PercentageComplete string `json:"percentageComplete"`
	TraceData          struct {
		SequenceFlow6 struct {
			TracePages struct {
				Num1 struct {
					Headers struct {
						SapFunctionName string `json:"sapFunctionName"`
					} `json:"headers"`
					Properties struct {
						CamelMessageHistory      string `json:"CamelMessageHistory"`
						CamelExternalRedelivered string `json:"CamelExternalRedelivered"`
						CamelToEndpoint          string `json:"CamelToEndpoint"`
						User                     string `json:"user"`
						Passwd                   string `json:"passwd"`
						Key                      string `json:"key"`
						Cert                     string `json:"cert"`
					} `json:"properties"`
					Body []int `json:"body"`
				} `json:"1"`
			} `json:"tracePages"`
			TraceState string `json:"traceState"`
		} `json:"SequenceFlow_6"`
		SequenceFlow5 struct {
			TracePages struct {
				Num1 struct {
					Headers struct {
					} `json:"headers"`
					Properties struct {
						CamelMessageHistory      string `json:"CamelMessageHistory"`
						CamelExternalRedelivered string `json:"CamelExternalRedelivered"`
						CamelToEndpoint          string `json:"CamelToEndpoint"`
						CamelCreatedTimestamp    string `json:"CamelCreatedTimestamp"`
					} `json:"properties"`
					Body []interface{} `json:"body"`
				} `json:"1"`
			} `json:"tracePages"`
			TraceState string `json:"traceState"`
		} `json:"SequenceFlow_5"`
	} `json:"traceData"`
}
