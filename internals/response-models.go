package internals

type AdmissionStatus struct {
	Message string `json:"message"`
}

type AdmissionResponse struct {
	UID       string           `json:"uid"`
	Allowed   bool             `json:"allowed"`
	Patch     string           `json:"patch,omitempty"`
	PatchType string           `json:"patchType,omitempty"`
	Status    *AdmissionStatus `json:"status,omitempty"`
}

type AdmissionReviewResponse struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Response   *AdmissionResponse `json:"response"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
