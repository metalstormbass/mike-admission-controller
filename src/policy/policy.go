package policy

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/open-policy-agent/opa/rego"

	"k8s.io/apimachinery/pkg/types"
)

// Basic Policy without using rego.
func ValidateContainerTag(image string) (valid bool) {

	imageTagMap := strings.SplitN(image, ":", 2)
	imageTag := string(imageTagMap[1])
	if imageTag != "latest" {
		return true
	} else {
		return false
	}

}

// Load Rego Policy

//

const opaPolicy = `
package admissionControl

deny[msg] {
    input.request.object.apiVersion == "apps/v1"
    input.request.object.kind == "Deployment"
    
    container := input.request.object.spec.template.spec.containers[_]
    endswith(container.image, ":latest")
    msg = sprintf("Container '%v'is using 'latest' tag", [container.name])
}

`

// Admission Control Function

func AdmissionController(resource interface{}, uid types.UID, w http.ResponseWriter) {
	ctx := context.Background()

	query, err := rego.New(
		rego.Query("data.admissionControl.deny"),
		rego.Module("policy.rego", opaPolicy),
	).PrepareForEval(ctx)

	if err != nil {
		log.Print(err)
		return
	}

	rs, err := query.Eval(ctx, rego.EvalInput(map[string]interface{}{
		"request": map[string]interface{}{
			"object": resource,
		},
	}))
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(rs.Allowed())
	log.Print(rs)
	if rs[0].Expressions[0].Value == nil {
		log.Print("Container was allowed")
	} else if rs[0].Expressions[0].Value != nil {
		log.Print("Container was NOT allowed")
	}

}

/*
	ctx := context.Background()

	query, err := rego.New(
		rego.Query("data.admissionControl.deny"),
		rego.Module("policy.rego", opaPolicy),
	).PrepareForEval(ctx)

	if err != nil {
		log.Print(err)
		return
	}

	rs, err := query.Eval(ctx, rego.EvalInput(map[string]interface{}{
		"request": map[string]interface{}{
			"object": resource,
		},
	}))

	if err != nil {
		log.Print(err)
		return
	}
	log.Print(rs[0].Expressions[0].Value.(string))
	/*if len(rs) > 0 && rs[0].Expressions[0].Value != nil {
		denyMsg := rs[0].Expressions[0].Value.(string)
		log.Printf("Admission Denied: %s\n", denyMsg)
		//return
	} else {
		log.Println("Admission Allowed")
		//return
	}
*/
