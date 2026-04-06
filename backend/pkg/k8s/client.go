package k8s

import (
	"context"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
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

// AnalyzeBlastRadius simulează/extrage permisiunile de rețea și RBAC
func AnalyzeBlastRadius(namespace string, podName string) (rbacRisk string, networkRisk string) {
	if Clientset == nil {
		return "Unknown", "Unknown"
	}

	pod, err := Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "Unknown", "Unknown"
	}

	// RBAC Logic
	sa := pod.Spec.ServiceAccountName
	if sa == "default" || sa == "" {
		rbacRisk = "Safe: Default SA"
	} else {
		rbacRisk = "Critical RBAC Privileges: Potential Secret/Admin access"
	}

	// Network Logic (Simplificat: presupunem că vrem să dăm avertisment dacă găsim o imagine web)
	networkRisk = "Egress: Unrestricted. Pod can communicate with the external internet."
	
	// Analizeaza pe viitor networking.k8s.io/v1 NetworkPolicies aici.

	return rbacRisk, networkRisk
}
