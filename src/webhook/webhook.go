package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/metalstormbass/mike-admission-controller/src/policy"
	"github.com/rs/zerolog/log"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Webhook

func Validate(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview

	// Parse Input

	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print("Error Decoding")
		return
	}

	// Extract UID

	uid := admissionReview.Request.UID

	// Extract Pods

	var pod corev1.Pod
	json.Unmarshal(admissionReview.Request.Object.Raw, &pod)

	// Extract Image String and Check for tag
	checkTag(pod, uid, w)

}

func checkTag(pod corev1.Pod, uid types.UID, w http.ResponseWriter) {
	var image string

	for i := range pod.Spec.Containers {
		image = pod.Spec.Containers[i].Image

		// Build Response Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var admissionResponse admissionv1.AdmissionResponse
		if policy.ValidateContainerTag(image) {

			// Define Response
			admissionResponse = admissionv1.AdmissionResponse{
				Allowed: true,
				UID:     uid,

				Result: &metav1.Status{
					Status:  metav1.StatusSuccess,
					Message: "Container tag validation succeeded",
					Code:    http.StatusOK,
				},
			}

		} else {

			// Define Respnse

			admissionResponse = admissionv1.AdmissionResponse{
				Allowed: false,
				UID:     uid,

				Result: &metav1.Status{
					Status:  metav1.StatusFailure,
					Message: "Container tag validation failed",
					Reason:  metav1.StatusReasonForbidden,
					Code:    http.StatusForbidden,
					Details: &metav1.StatusDetails{
						Causes: []metav1.StatusCause{
							{
								Message: "Invalid container tag",
								Field:   "spec.containers[" + string(i) + "].image",
							},
						},
					},
				},
			}
		}

		admissionReviewResponse := admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "admission.k8s.io/v1",
				Kind:       "AdmissionReview",
			},
			Response: &admissionResponse,
		}
		log.Print(admissionReviewResponse)
		err := json.NewEncoder(w).Encode(admissionReviewResponse)
		if err != nil {
			http.Error(w, "Failed to encode AdmissionReview response", http.StatusInternalServerError)
			return
		}

	}
}

/*
Just keeping this here for reference

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

	   	if policy.ValidateContainerTag(image) {
	   		log.Print("This is allowed")
	   	} else {
	   		log.Print("NOT allowed")
	   	}

	   }
*/
