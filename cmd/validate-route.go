package main

import (
	"encoding/json"
	"net/http"
	"simple-admission-webhook/internals"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var admissionReview internals.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		json.NewEncoder(w).Encode(internals.AdmissionResponse{
			Allowed: false,
			Status:  internals.AdmissionStatus{Message: "Invalid JSON: " + err.Error()},
		})
		return
	}

	objectName := admissionReview.Request.Object.Metadata.Name
	userName := admissionReview.Request.UserInfo.Username

	message := "Validations successful"
	status := true
	if objectName == userName {
		message = "Validations failed: object name cannot be the same as username"
		status = false
	}

	var statusObj internals.AdmissionStatus = internals.AdmissionStatus{
		Message: message,
	}

	json.NewEncoder(w).Encode(internals.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: status,
		Status:  statusObj,
	})
}
