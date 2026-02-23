package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"simple-admission-webhook/internals"
)

func Mutate(w http.ResponseWriter, r *http.Request) {
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

	objectContainers := admissionReview.Request.Object.Spec.Containers
	if len(objectContainers) == 0 {
		println("Failed to parse containers from request")
		json.NewEncoder(w).Encode(internals.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: false,
			Status:  internals.AdmissionStatus{Message: "Failed to parse containers from request"},
		})
		return
	}

	securityContext := objectContainers[0].SecurityContext
	var patches []internals.PatchOperation
	if securityContext == nil {
		patches = append(patches, internals.PatchOperation{
			Op:   "add",
			Path: "/spec/containers/0/securityContext",
			Value: map[string]int{
				"runAsUser": 1000,
			},
		})
	} else if securityContext.RunAsUser == nil || *securityContext.RunAsUser != 1000 {
		patches = append(patches, internals.PatchOperation{
			Op:    "add",
			Path:  "/spec/containers/0/securityContext/runAsUser",
			Value: 1000,
		})
	}

	patchString := ""
	patchType := ""
	allowed := true
	if len(patches) > 0 {
		patchBytes, err := json.Marshal(patches)
		if err != nil {
			allowed = false
		} else {
			patchString = base64.StdEncoding.EncodeToString(patchBytes)
			patchType = "JSONPatch"
		}
	}

	response := internals.AdmissionMutateResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   allowed,
		Patch:     patchString,
		PatchType: patchType,
	}

	json.NewEncoder(w).Encode(response)
}
