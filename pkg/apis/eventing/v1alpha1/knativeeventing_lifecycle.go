package v1alpha1

import (
	"github.com/knative/pkg/apis"
)

var conditions = apis.NewLivingConditionSet(
	DeploymentsAvailable,
	InstallSucceeded,
)

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}

func (s *KnativeEventingStatus) IsReady() bool {
	return conditions.Manage(s).IsHappy()
}

func (s *KnativeEventingStatus) IsInstalled() bool {
	return s.GetCondition(InstallSucceeded).IsTrue()
}

func (s *KnativeEventingStatus) IsAvailable() bool {
	return s.GetCondition(DeploymentsAvailable).IsTrue()
}

func (s *KnativeEventingStatus) IsDeploying() bool {
	return s.IsInstalled() && !s.IsAvailable()
}

func (s *KnativeEventingStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return conditions.Manage(s).GetCondition(t)
}

func (s *KnativeEventingStatus) InitializeConditions() {
	conditions.Manage(s).InitializeConditions()
}

func (s *KnativeEventingStatus) MarkInstallFailed(msg string) {
	conditions.Manage(s).MarkFalse(
		InstallSucceeded,
		"Error",
		"Install failed with message: %s", msg)
}

func (s *KnativeEventingStatus) MarkInstallSucceeded() {
	conditions.Manage(s).MarkTrue(InstallSucceeded)
}

func (s *KnativeEventingStatus) MarkDeploymentsAvailable() {
	conditions.Manage(s).MarkTrue(DeploymentsAvailable)
}

func (s *KnativeEventingStatus) MarkDeploymentsNotReady() {
	conditions.Manage(s).MarkFalse(
		DeploymentsAvailable,
		"NotReady",
		"Waiting on deployments")
}

func (s *KnativeEventingStatus) MarkIgnored(msg string) {
	conditions.Manage(s).MarkFalse(
		InstallSucceeded,
		"Ignored",
		"Install not attempted: %s", msg)
}
