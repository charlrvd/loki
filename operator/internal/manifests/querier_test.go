package manifests_test

import (
	"testing"

	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
	"github.com/grafana/loki/operator/internal/manifests"
	"github.com/stretchr/testify/require"
)

func TestNewQuerierDeployment_HasTemplateConfigHashAnnotation(t *testing.T) {
	ss := manifests.NewQuerierDeployment(manifests.Options{
		Name:       "abcd",
		Namespace:  "efgh",
		ConfigSHA1: "deadbeef",
		Stack: lokiv1.LokiStackSpec{
			StorageClassName: "standard",
			Template: &lokiv1.LokiTemplateSpec{
				Querier: &lokiv1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	expected := "loki.grafana.com/config-hash"
	annotations := ss.Spec.Template.Annotations
	require.Contains(t, annotations, expected)
	require.Equal(t, annotations[expected], "deadbeef")
}

func TestNewQuerierDeployment_HasTemplateCertRotationRequiredAtAnnotation(t *testing.T) {
	ss := manifests.NewQuerierDeployment(manifests.Options{
		Name:                   "abcd",
		Namespace:              "efgh",
		CertRotationRequiredAt: "deadbeef",
		Stack: lokiv1.LokiStackSpec{
			Template: &lokiv1.LokiTemplateSpec{
				Querier: &lokiv1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	expected := "loki.grafana.com/certRotationRequiredAt"
	annotations := ss.Spec.Template.Annotations
	require.Contains(t, annotations, expected)
	require.Equal(t, annotations[expected], "deadbeef")
}

func TestNewQuerierDeployment_SelectorMatchesLabels(t *testing.T) {
	// You must set the .spec.selector field of a Deployment to match the labels of
	// its .spec.template.metadata.labels. Prior to Kubernetes 1.8, the
	// .spec.selector field was defaulted when omitted. In 1.8 and later versions,
	// failing to specify a matching Pod Selector will result in a validation error
	// during Deployment creation.
	// See https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#pod-selector
	ss := manifests.NewQuerierDeployment(manifests.Options{
		Name:      "abcd",
		Namespace: "efgh",
		Stack: lokiv1.LokiStackSpec{
			StorageClassName: "standard",
			Template: &lokiv1.LokiTemplateSpec{
				Querier: &lokiv1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	l := ss.Spec.Template.GetObjectMeta().GetLabels()
	for key, value := range ss.Spec.Selector.MatchLabels {
		require.Contains(t, l, key)
		require.Equal(t, l[key], value)
	}
}
