package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/metalstormbass/mike-admission-controller/src/policy"
	"github.com/rs/zerolog/log"
	admission "k8s.io/api/admission/v1"
	core "k8s.io/api/core/v1"
)

// Response Struct

func Validate(w http.ResponseWriter, r *http.Request) {
	var admissionReview admission.AdmissionReview

	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print("Error Decoding")
		return
	}
	var pod core.Pod
	json.Unmarshal(admissionReview.Request.Object.Raw, &pod)

	var image string
	for i := range pod.Spec.Containers {
		image = pod.Spec.Containers[i].Image
		if policy.ValidateContainerTag(image) {

			// Return the AdmissionResponse
			admissionResponse := admission.AdmissionResponse{
				Allowed: true,
			}

			// Encode and send the AdmissionReview response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(admissionResponse)
			if err != nil {
				http.Error(w, "Failed to encode AdmissionReview response", http.StatusInternalServerError)
				return
			}

		} else {
			log.Print("NOT allowed")
		}
	}

	/*

		var data interface{}

		// Convert to a blob of JSON
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Make JSON Parseable
		jsonObj, ok := data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid JSON object", http.StatusBadRequest)
			return
		}

		containers := jsonObj["request"].(map[string]interface{})["object"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})

		// Iterate over the containers
		var image string
		for _, container := range containers {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				log.Error()
			}
			image = containerMap["image"].(string)
		}

	*/
	// Define Response

	/*
	   	if policy.ValidateContainerTag(image) {
	   		log.Print("This is allowed")
	   	} else {
	   		log.Print("NOT allowed")
	   	}

	   }
	*/
}
