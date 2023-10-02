/*
Copyright 2022-2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Constants used with SnapshotEnvironmentBindingStatus's ComponentDeploymentConditions field
const (
	ComponentDeploymentConditionAllComponentsDeployed = "AllComponentsDeployed"
	ComponentDeploymentConditionCommitsSynced         = "CommitsSynced"
	ComponentDeploymentConditionCommitsUnsynced       = "CommitsUnsynced"
	ComponentDeploymentConditionErrorOccurred         = "ErrorOccurred"
)

// See 'SnapshotEnvironmentBinding' resource for details of this resource. SnapshotEnvironmentBindingSpec defines the desired state of SnapshotEnvironmentBinding.
type SnapshotEnvironmentBindingSpec struct {

	// Application is a reference to the Application resource (defined in the same namespace) that we are deploying as part of this SnapshotEnvironmentBinding.
	// Required
	// +required
	Application string `json:"application"`

	// Environment is the environment resource (defined in the namespace) that the binding will deploy to.
	// Required
	// +required
	Environment string `json:"environment"`

	// Snapshot is the Snapshot resource (defined in the namespace) that contains the container image versions
	// for the components of the Application.
	// Required
	// +required
	Snapshot string `json:"snapshot"`

	// Component-specific configuration information, used when generating GitOps repository resources.
	// Required.
	// +required
	Components []BindingComponent `json:"components"`
}

// BindingComponent contains individual component data
type BindingComponent struct {

	// Name is the name of the component.
	Name string `json:"name"`

	// Configuration describes GitOps repository customizations that are specific to the
	// the component-application-environment combination.
	// - Values defined in this struct will overwrite values from Application/Environment/Component.
	// Optional
	// +optional
	Configuration BindingComponentConfiguration `json:"configuration,omitempty"`
}

// BindingComponentConfiguration describes GitOps repository customizations that are specific to the
// the component-application-environment combination.
type BindingComponentConfiguration struct {
	// NOTE: The specific fields, and their form, to be included are TBD.

	// API discussion concluded with no obvious need for target port; it is thus excluded here.
	// Let us know if you have a requirement here.
	//
	// TargetPort int `json:"targetPort"`

	// Replicas defines the number of replicas to use for the component
	// Optional
	// +optional
	Replicas *int `json:"replicas,omitempty"`

	// Resources defines the Compute Resources required by the component.
	// Optional.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Env describes environment variables to use for the component.
	// Optional.
	// +optional
	Env []EnvVarPair `json:"env,omitempty"`
}

// EnvVarPair describes environment variables to use for the component
type EnvVarPair struct {

	// Name is the environment variable name
	Name string `json:"name"`

	// Value is the environment variable value
	Value string `json:"value"`
}

// BindingComponentGitOpsRepository is a reference to a GitOps repository, including path/branch
// where the application/component/environment resources can be found (usually via a kustomize overlay).
type BindingComponentGitOpsRepository struct {

	// URL is the Git repository URL
	// e.g. The Git repository that contains the K8s resources to deployment for the component of the application.
	URL string `json:"url"`

	// Branch is the branch to use when accessing the GitOps repository
	Branch string `json:"branch"`

	// Path is a pointer to a folder in the GitOps repo, containing a kustomization.yaml
	// NOTE: Each component-env combination must have it's own separate path
	Path string `json:"path"`

	// GeneratedResources contains the list of GitOps repository resources generated by the application service controller
	// in the overlays/<environment> dir, for example, 'deployment-patch.yaml'. This is stored to differentiate between
	// application-service controller generated resources vs resources added by a user
	GeneratedResources []string `json:"generatedResources"`

	// CommitID contains the most recent commit ID for which the Kubernetes resources of the Component were modified.
	CommitID string `json:"commitID"`
}

// BindingComponentStatus contains the status of the components
type BindingComponentStatus struct {

	// Name is the name of the component.
	Name string `json:"name"`

	// GeneratedRouteName is the name of the route that was generated for the Component, if a Route was generated.
	GeneratedRouteName string `json:"generatedRouteName,omitempty"`

	// GitOpsRepository contains the Git URL, path, branch, and most recent commit id for the component
	GitOpsRepository BindingComponentGitOpsRepository `json:"gitopsRepository"`
}

// SnapshotEnvironmentBindingStatus defines the observed state of SnapshotEnvironmentBinding
type SnapshotEnvironmentBindingStatus struct {

	// GitOpsDeployments describes the set of GitOpsDeployment resources that are owned by the SnapshotEnvironmentBinding, and are
	// deploying the Components of the Application to the target Environment.
	// To determine the health/sync status of a binding, you can look at the GitOpsDeployments decribed here.
	GitOpsDeployments []BindingStatusGitOpsDeployment `json:"gitopsDeployments,omitempty"`

	// Components describes a component's GitOps repository information.
	// This status is updated by the Application Service controller.
	Components []BindingComponentStatus `json:"components,omitempty"`

	// Condition describes operations on the GitOps repository, for example, if there were issues with generating/processing the repository.
	// This status is updated by the Application Service controller.
	GitOpsRepoConditions []metav1.Condition `json:"gitopsRepoConditions,omitempty"`

	// BindingConditions will contain user-oriented error messages from the SnapshotEnvironmentBinding reconciler.
	BindingConditions []metav1.Condition `json:"bindingConditions,omitempty"`

	// ComponentDeploymentConditions describes the deployment status of all of the Components of the Application.
	// This status is updated by the Gitops Service's SnapshotEnvironmentBinding controller
	ComponentDeploymentConditions []metav1.Condition `json:"componentDeploymentConditions,omitempty"`
}

// BindingStatusGitOpsDeployment describes an individual reference
// to a GitOpsDeployment resources that is used to deploy this binding.
//
// To determine the health/sync status of a binding, you can look at the GitOpsDeployments decribed here.
type BindingStatusGitOpsDeployment struct {

	// ComponentName is the name of the component in the (component, gitopsdeployment) pair
	ComponentName string `json:"componentName"`

	// GitOpsDeployment is a reference to the name of a GitOpsDeployment resource which is used to deploy the binding.
	// The Health/sync status for the binding can thus be read from the references GitOpsDeployment
	GitOpsDeployment string `json:"gitopsDeployment,omitempty"`

	// GitOpsDeploymentSyncStatus is the sync status of the deployment owned by the binding
	GitOpsDeploymentSyncStatus string `json:"syncStatus,omitempty"`

	// GitOpsDeploymentHealthStatus is the health status of the deployment owned by the binding
	GitOpsDeploymentHealthStatus string `json:"health,omitempty"`

	// GitOpsDeploymentCommitID is the commit ID of the GitOpsDeployment
	GitOpsDeploymentCommitID string `json:"commitID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// The `SnapshotEnvironmentBinding` resource specifies the deployment relationship between (a single application, a
// single environment, and a single snapshot) combination.
//
// It can be thought of as a 3-tuple that defines what Application should be deployed to what Environment, and
// which Snapshot should be deployed (Snapshot being the specific component container image versions of that
// Aplication that should be deployed to that Environment).
//
// **Note**: There should not exist multiple SnapshotEnvironmentBinding CRs in a Namespace that share the same
// Application and Environment value. For example:
// - Good:
//   - SnapshotEnvironmentBinding A: (application=appA, environment=dev, snapshot=my-snapshot)
//   - SnapshotEnvironmentBinding B: (application=appA, environment=staging, snapshot=my-snapshot)
//
// - Bad:
//   - SnapshotEnvironmentBinding A: (application=*appA*, environment=*staging*, snapshot=my-snapshot)
//   - SnapshotEnvironmentBinding B: (application=*appA*, environment=*staging*, snapshot=second-snapshot)
//
// +kubebuilder:resource:path=snapshotenvironmentbindings,shortName=aseb;binding
type SnapshotEnvironmentBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotEnvironmentBindingSpec   `json:"spec"`
	Status SnapshotEnvironmentBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SnapshotEnvironmentBindingList contains a list of SnapshotEnvironmentBinding
type SnapshotEnvironmentBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnapshotEnvironmentBinding `json:"items"`
}