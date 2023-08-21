package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/metalstormbass/mike-admission-controller/src/policy"
	"github.com/rs/zerolog/log"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
)

// Define a common interface for different resource types
type KubernetesResource interface {
	GetKind() string
}

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
	kind := (admissionReview.Request.Kind.Kind)

	// Depending on the kind, unmarshal into appropriate struct
	var resource interface{}

	switch kind {
	case "Deployment":
		var deployment appsv1.Deployment
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &deployment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print("Error unmarshaling Deployment")
			return
		}
		resource = &deployment
	case "Pod":
		var pod corev1.Pod
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print("Error unmarshaling Pod")
			return
		}
		resource = &pod
	case "Service":
		var service corev1.Service
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print("Error unmarshaling Service")
			return
		}
		resource = &service
	default:
		log.Printf("%s is not a supported type", kind)
		return
	}
	// Extract Image String and Check for tag

	policy.AdmissionController(resource, uid, w)

}

/*

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
*/
