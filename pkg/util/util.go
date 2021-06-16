package util

import (
	"cf-argo-plugin/pkg/codefresh"
	"errors"
	"fmt"
)

func FilterActivity(applicationName string, updatedActivities []codefresh.UpdatedActivity) (error, codefresh.UpdatedActivity) {
	var rolloutActivity codefresh.UpdatedActivity
	for _, activity := range updatedActivities {
		if activity.ApplicationName == applicationName {
			return nil, activity
		}
	}
	return errors.New(fmt.Sprintf("can't find activity with app name %s", applicationName)), rolloutActivity
}
