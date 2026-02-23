package internals

type PodSpec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name            string           `json:"name"`
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
}

type SecurityContext struct {
	RunAsUser *int64 `json:"runAsUser,omitempty"`
}

type AdmissionReview struct {
	Request struct {
		UID      string `json:"uid"`
		UserInfo struct {
			Username string `json:"username"`
		} `json:"userInfo"`

		Object struct {
			Metadata struct {
				Name string `json:"name"`
				Kind string `json:"kind"`
			} `json:"metadata"`
			Spec PodSpec `json:"spec"`
		} `json:"object"`
	} `json:"request"`
}
