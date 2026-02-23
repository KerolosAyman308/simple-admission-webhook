package internals

type AdmissionStatus struct {
	Message string `json:"message"`
}

type AdmissionResponse struct {
	UID     string          `json:"uid"`
	Allowed bool            `json:"allowed"`
	Status  AdmissionStatus `json:"status"`
}

type AdmissionMutateResponse struct {
	UID       string `json:"uid"`
	Allowed   bool   `json:"allowed"`
	Patch     string `json:"patch"`
	PatchType string `json:"patchType"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
