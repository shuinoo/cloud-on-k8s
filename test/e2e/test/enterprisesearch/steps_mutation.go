// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package enterprisesearch

import (
	"testing"

	entv1beta1 "github.com/elastic/cloud-on-k8s/pkg/apis/enterprisesearch/v1beta1"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/hash"
	"github.com/elastic/cloud-on-k8s/pkg/utils/k8s"
	"github.com/elastic/cloud-on-k8s/test/e2e/test"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func (b Builder) MutationTestSteps(k *test.K8sClient) test.StepList {
	return b.AnnotatePodsWithBuilderHash(k).
		WithSteps(b.UpgradeTestSteps(k)).
		WithSteps(b.CheckK8sTestSteps(k)).
		WithSteps(b.CheckStackTestSteps(k))
}

func (b Builder) AnnotatePodsWithBuilderHash(k *test.K8sClient) test.StepList {
	return []test.Step{
		{
			Name: "Annotate Pods with a hash of their Builder spec",
			Test: test.Eventually(func() error {
				var pods corev1.PodList
				if err := k.Client.List(&pods, test.EnterpriseSearchPodListOptions(b.EnterpriseSearch.Namespace, b.EnterpriseSearch.Name)...); err != nil {
					return err
				}

				expectedHash := hash.HashObject(b.MutatedFrom.EnterpriseSearch.Spec)
				for _, pod := range pods.Items {
					if err := test.AnnotatePodWithBuilderHash(k, pod, expectedHash); err != nil {
						return err
					}
				}
				return nil
			}),
		},
	}
}

func (b Builder) MutationReversalTestContext() test.ReversalTestContext {
	panic("not implemented")
}

func (b Builder) UpgradeTestSteps(k *test.K8sClient) test.StepList {
	return test.StepList{
		{
			Name: "Applying the Enterprise Search mutation should succeed",
			Test: func(t *testing.T) {
				var ent entv1beta1.EnterpriseSearch
				require.NoError(t, k.Client.Get(k8s.ExtractNamespacedName(&b.EnterpriseSearch), &ent))
				ent.Spec = b.EnterpriseSearch.Spec
				require.NoError(t, k.Client.Update(&ent))
			},
		}}
}
