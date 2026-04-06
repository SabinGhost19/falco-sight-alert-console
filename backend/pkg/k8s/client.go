package k8s

import (
"context"
"log"
"strconv"
"strings"

"gopkg.in/yaml.v3"
corev1 "k8s.io/api/core/v1"
netv1 "k8s.io/api/networking/v1"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/rest"
"k8s.io/client-go/tools/clientcmd"
)

var Clientset *kubernetes.Clientset

// InitK8sClient folosește kubeconfig din sistem (sau in-cluster)
func InitK8sClient(kubeconfigPath string) {
	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		// Fallback for in-cluster configuration
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %v", err)
	}
	log.Println("Kubernetes client initialized successfully.")
}

// FetchPodManifest extrage reprezentarea YAML live a unui Pod
func FetchPodManifest(namespace string, podName string) (string, error) {
	pod, err := Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Sterge datele de status si metadata inutile pentru claritate
	pod.Status = corev1.PodStatus{} // reset
	pod.ObjectMeta.ManagedFields = nil

	yamlBytes, err := yaml.Marshal(pod)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// AnalyzeManifest realizează o evaluare statică simplă a manifestului
func AnalyzeManifest(manifestYAML string) string {
	var vulnerableLines []string

	// Analizează vizual YAML-ul pentru security risks
	lines := strings.Split(manifestYAML, "\n")
	for i, line := range lines {
		// O heuristică simplă: dacă există `privileged: true`
		if strings.Contains(line, "privileged: true") || strings.Contains(line, "runAsRoot: true") {
			// +1 pentru că editorul (ex: Monaco) e indexat 1-based
			vulnerableLines = append(vulnerableLines, strconv.Itoa(i+1))
		}
	}
	return strings.Join(vulnerableLines, ",")
}

// AnalyzeBlastRadius simulează/extrage permisiunile de rețea și RBAC real din Kubernetes
func AnalyzeBlastRadius(namespace string, podName string) (rbacRisk string, networkRisk string) {
	if Clientset == nil {
		return "Unknown: k8s client disabled", "Unknown"
	}

	pod, err := Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "Unknown: pod missing", "Unknown"
	}

	// 1. RBAC Extraction
	sa := pod.Spec.ServiceAccountName
	if sa == "" {
		sa = "default"
	}

	rbacRisk = "Safe: Default SA"
	if sa != "default" {
		rbacRisk = "Elevated: " + sa + " Service Account."

		// Attempt to parse deep bindings 
		// (Simplificat aici in mod real se pot itera ClusterRoleBindings / RoleBindings, dar pentru scop demo adaugam info detaliat)
		crbs, err := Clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			for _, crb := range crbs.Items {
				for _, subject := range crb.Subjects {
					if subject.Kind == "ServiceAccount" && subject.Name == sa && subject.Namespace == namespace {
						rbacRisk = "Critical RBAC Privileges: Pod mapped to ClusterRole '" + crb.RoleRef.Name + "'"
						break
					}
				}
			}
		}
	}

	// 2. Network Policies Extraction
	netPol, err := Clientset.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil || len(netPol.Items) == 0 {
		networkRisk = "Egress: Unrestricted. No NetworkPolicy found. Pod can connect anywhere."
	} else {
		// Daca avem elemente - cautam daca vreuna influenteaza ingress/egress ul podului...
		networkRisk = "Protected: Network Policies are defined in this namespace."
		hasEgressLimit := false
		for _, np := range netPol.Items {
			for _, p := range np.Spec.PolicyTypes {
				if p == netv1.PolicyTypeEgress {
					hasEgressLimit = true
				}
			}
		}
		if !hasEgressLimit {
			networkRisk = "Partial Risk: Policies exist but Egress is unrestricted. Exfiltration possible."
		}
	}

	return rbacRisk, networkRisk
}
