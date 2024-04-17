package exporters

import (
	"context"
	"fmt"

	"github.com/redhat-appstudio/rhtap-cli/cmd/rhtap/commands/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TransformSecret(ctx context.Context, fetchedResourceList runtime.Object, cloneConfig config.CloneConfig, localCache []runtime.Object) ([]runtime.Object, error) {

	var selectedResources []runtime.Object

	itsList, ok := fetchedResourceList.(*v1.SecretList)

	if !ok {
		return nil, fmt.Errorf("resource of type SecretList was not passsed")
	}

	for _, secret := range itsList.Items {
		if secret.Name != "snyk-secret" {
			continue
		}

		// TODO: Encrypt the content of the secret
		copiedSecret := secret
		selectedResources = append(selectedResources, &copiedSecret)
	}
	return selectedResources, nil
}

func GenerateYAMLForSecret(ctx context.Context, transformedResources []runtime.Object) ([][]byte, error) {
	var resourcesInYAML [][]byte
	for _, resource := range transformedResources {
		secret := resource.(*v1.Secret)
		inBytes, err := yaml.Marshal(secret)
		if err != nil {
			return nil, err
		}
		resourcesInYAML = append(resourcesInYAML, inBytes)
	}
	return resourcesInYAML, nil
}

func GetSecrets(ctx context.Context, namespace string, cloneConfig config.CloneConfig, client *kubernetes.Clientset) (runtime.Object, error) {
	var secrets = &v1.SecretList{}
	client.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/appstudio.redhat.com/v1beta1/namespaces/%s/secrets/snyk-secret", namespace)).
		Do(context.TODO()).Into(secrets)

	//TODO: Return error but handle/ignore it was due to lack of permissions.
	//return secrets, err
	return secrets, nil
}
