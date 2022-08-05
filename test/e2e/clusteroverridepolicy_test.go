package e2e

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog/v2"

	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	"github.com/karmada-io/karmada/test/e2e/framework"
	testhelper "github.com/karmada-io/karmada/test/helper"
)

var _ = ginkgo.Describe("[BasicClusterOverridePolicy] basic cluster override policy testing", func() {
	ginkgo.Context("Namespace propagation testing", func() {
		var (
			randNamespace string
			ns            *corev1.Namespace
			namespaceCop  *policyv1alpha1.ClusterOverridePolicy

			customLabelKey          = "hello"
			customLabelVal          = "world"
			plaintextOverriderValue = fmt.Sprintf("{ \"%s\": \"%s\"}", customLabelKey, customLabelVal)
		)

		ginkgo.BeforeEach(func() {
			randNamespace = "cop-test-ns-" + rand.String(RandomStrLength)
			ns = testhelper.NewNamespace(randNamespace)
			namespaceCop = testhelper.NewClusterOverridePolicyByOverrideRules(ns.Name,
				[]policyv1alpha1.ResourceSelector{
					{
						APIVersion: "v1",
						Kind:       "Namespace",
						Name:       ns.Name,
					},
				},
				[]policyv1alpha1.RuleWithCluster{
					{
						TargetCluster: &policyv1alpha1.ClusterAffinity{
							ClusterNames: framework.ClusterNames(),
						},
						Overriders: policyv1alpha1.Overriders{
							Plaintext: []policyv1alpha1.PlaintextOverrider{
								{
									Path:     "/metadata/labels",
									Operator: policyv1alpha1.OverriderOpAdd,
									Value:    apiextensionsv1.JSON{Raw: []byte(plaintextOverriderValue)},
								},
							},
						},
					},
				})
		})

		ginkgo.BeforeEach(func() {
			framework.CreateClusterOverridePolicy(karmadaClient, namespaceCop)
			framework.CreateNamespace(kubeClient, ns)
			ginkgo.DeferCleanup(func() {
				framework.RemoveClusterOverridePolicy(karmadaClient, namespaceCop.Name)
				framework.RemoveNamespace(kubeClient, ns.Name)
				framework.WaitNamespaceDisappearOnClusters(framework.ClusterNames(), ns.Name)
			})
		})

		ginkgo.It("Namespace testing", func() {
			ginkgo.By(fmt.Sprintf("Check if namespace(%s) present on member clusters", ns.Name), func() {
				for _, clusterName := range framework.ClusterNames() {
					clusterClient := framework.GetClusterClient(clusterName)
					gomega.Expect(clusterClient).ShouldNot(gomega.BeNil())

					klog.Infof("Waiting for namespace present on cluster(%s)", clusterName)
					gomega.Eventually(func(g gomega.Gomega) (bool, error) {
						clusterNs, err := clusterClient.CoreV1().Namespaces().Get(context.TODO(), ns.Name, metav1.GetOptions{})
						if err != nil {
							if apierrors.IsNotFound(err) {
								return false, nil
							}

							return false, err
						}

						v, ok := clusterNs.Labels[customLabelKey]
						if ok && v == customLabelVal {
							return true, nil
						}
						return false, nil
					}, pollTimeout, pollInterval).Should(gomega.Equal(true))
				}
			})
		})
	})
})
