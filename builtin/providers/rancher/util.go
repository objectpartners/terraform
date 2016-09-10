package rancher

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s", value, k))
		}
		return
	}
}

type getState func() string

type getType func(id string) (getState, error)

func waitForStatus(state, id string, getter getType) error {
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		stateGetter, e := getter(id)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		currentState := stateGetter()
		if currentState != state {
			return resource.RetryableError(fmt.Errorf("Object[%s] is not in state %s[%s].", id, state, currentState))
		}
		return nil
	})
	return err
}
