package aws

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Data types returned to the UI layer
// ---------------------------------------------------------------------------

// K8sNamespace represents a Kubernetes namespace.
type K8sNamespace struct {
	Name   string
	Status string
}

// K8sPod represents a Kubernetes pod.
type K8sPod struct {
	Name       string
	Namespace  string
	Node       string
	Status     string
	IP         string
	Ready      string // "2/2"
	Restarts   int
	Age        string
	Containers []K8sPodContainer
}

// K8sPodContainer is a container within a pod.
type K8sPodContainer struct {
	Name         string
	Image        string
	State        string
	Ready        bool
	RestartCount int
}

// K8sDeployment represents a Kubernetes deployment.
type K8sDeployment struct {
	Name      string
	Namespace string
	Replicas  string // "3/3"
	Ready     int
	UpToDate  int
	Available int
	Age       string
}

// K8sService represents a Kubernetes service.
type K8sService struct {
	Name       string
	Namespace  string
	Type       string // ClusterIP, NodePort, LoadBalancer
	ClusterIP  string
	ExternalIP string
	Ports      string // "80/TCP,443/TCP"
	Age        string
}

// K8sNode represents a Kubernetes node.
type K8sNode struct {
	Name       string
	Status     string
	Version    string
	Roles      string
	InternalIP string
	ExternalIP string
	OS         string
	Arch       string
	CPU        string
	Memory     string
	Age        string
}

// ---------------------------------------------------------------------------
// Token management
// ---------------------------------------------------------------------------

type cachedToken struct {
	token   string
	expires time.Time
}

var (
	tokenCache   = make(map[string]*cachedToken)
	tokenCacheMu sync.Mutex
)

// GetEKSToken obtains a bearer token for K8s API authentication via
// `aws eks get-token`. The token is cached for 14 minutes (expires at 15).
func GetEKSToken(ctx context.Context, clusterName, profile, region string) (string, error) {
	key := clusterName + "::" + profile + "::" + region
	tokenCacheMu.Lock()
	if ct, ok := tokenCache[key]; ok && time.Now().Before(ct.expires) {
		tokenCacheMu.Unlock()
		return ct.token, nil
	}
	tokenCacheMu.Unlock()

	args := []string{"eks", "get-token", "--cluster-name", clusterName, "--region", region, "--output", "json"}
	if profile != "" && profile != "default" && profile != InstanceRoleProfile {
		args = append(args, "--profile", profile)
	}

	cmd := exec.CommandContext(ctx, "aws", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("aws eks get-token: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("aws eks get-token: %w", err)
	}

	var resp struct {
		Status struct {
			Token string `json:"token"`
		} `json:"status"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}
	if resp.Status.Token == "" {
		return "", fmt.Errorf("empty token from aws eks get-token")
	}

	tokenCacheMu.Lock()
	tokenCache[key] = &cachedToken{token: resp.Status.Token, expires: time.Now().Add(14 * time.Minute)}
	tokenCacheMu.Unlock()

	return resp.Status.Token, nil
}

// ---------------------------------------------------------------------------
// K8s REST client
// ---------------------------------------------------------------------------

// K8sClient talks to the Kubernetes API of an EKS cluster.
type K8sClient struct {
	endpoint   string
	token      string
	httpClient *http.Client
}

// NewK8sClient creates a K8s REST client configured with the cluster CA and
// bearer token.
func NewK8sClient(endpoint, token, caCertBase64 string) (*K8sClient, error) {
	caData, err := base64.StdEncoding.DecodeString(caCertBase64)
	if err != nil {
		return nil, fmt.Errorf("decoding CA cert: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caData) {
		return nil, fmt.Errorf("invalid CA certificate")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	return &K8sClient{
		endpoint: strings.TrimRight(endpoint, "/"),
		token:    token,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		},
	}, nil
}

func (c *K8sClient) get(ctx context.Context, path string) ([]byte, error) {
	url := c.endpoint + path
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to extract K8s API error message
		var apiErr struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("K8s API %d: %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("K8s API %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ---------------------------------------------------------------------------
// Minimal JSON structs for K8s API responses
// ---------------------------------------------------------------------------

// We define only the fields we need from the K8s API.

type k8sListMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
}

// --- Namespace ---

type k8sNamespaceItem struct {
	Metadata k8sListMeta `json:"metadata"`
	Status   struct {
		Phase string `json:"phase"`
	} `json:"status"`
}

type k8sNamespaceList struct {
	Items []k8sNamespaceItem `json:"items"`
}

// --- Pod ---

type k8sPodItem struct {
	Metadata k8sListMeta `json:"metadata"`
	Spec     struct {
		NodeName   string `json:"nodeName"`
		Containers []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
		} `json:"containers"`
	} `json:"spec"`
	Status struct {
		Phase      string `json:"phase"`
		PodIP      string `json:"podIP"`
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"conditions"`
		ContainerStatuses []struct {
			Name         string `json:"name"`
			Image        string `json:"image"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			State        struct {
				Running    *struct{} `json:"running"`
				Waiting    *struct{ Reason string `json:"reason"` } `json:"waiting"`
				Terminated *struct{ Reason string `json:"reason"` } `json:"terminated"`
			} `json:"state"`
		} `json:"containerStatuses"`
		InitContainerStatuses []struct {
			Name         string `json:"name"`
			Image        string `json:"image"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			State        struct {
				Running    *struct{} `json:"running"`
				Waiting    *struct{ Reason string `json:"reason"` } `json:"waiting"`
				Terminated *struct{ Reason string `json:"reason"` } `json:"terminated"`
			} `json:"state"`
		} `json:"initContainerStatuses"`
	} `json:"status"`
}

type k8sPodList struct {
	Items []k8sPodItem `json:"items"`
}

// --- Deployment ---

type k8sDeploymentItem struct {
	Metadata k8sListMeta `json:"metadata"`
	Spec     struct {
		Replicas *int `json:"replicas"`
	} `json:"spec"`
	Status struct {
		Replicas          int `json:"replicas"`
		ReadyReplicas     int `json:"readyReplicas"`
		UpdatedReplicas   int `json:"updatedReplicas"`
		AvailableReplicas int `json:"availableReplicas"`
	} `json:"status"`
}

type k8sDeploymentList struct {
	Items []k8sDeploymentItem `json:"items"`
}

// --- Service ---

type k8sServiceItem struct {
	Metadata k8sListMeta `json:"metadata"`
	Spec     struct {
		Type      string `json:"type"`
		ClusterIP string `json:"clusterIP"`
		Ports     []struct {
			Port     int    `json:"port"`
			Protocol string `json:"protocol"`
			NodePort int    `json:"nodePort"`
		} `json:"ports"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
			Ingress []struct {
				Hostname string `json:"hostname"`
				IP       string `json:"ip"`
			} `json:"ingress"`
		} `json:"loadBalancer"`
	} `json:"status"`
}

type k8sServiceList struct {
	Items []k8sServiceItem `json:"items"`
}

// --- Node ---

type k8sNodeItem struct {
	Metadata k8sListMeta `json:"metadata"`
	Spec     struct{} `json:"spec"`
	Status   struct {
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"conditions"`
		Addresses []struct {
			Type    string `json:"type"`
			Address string `json:"address"`
		} `json:"addresses"`
		NodeInfo struct {
			KubeletVersion          string `json:"kubeletVersion"`
			OperatingSystem         string `json:"operatingSystem"`
			Architecture            string `json:"architecture"`
		} `json:"nodeInfo"`
		Allocatable map[string]string `json:"allocatable"`
	} `json:"status"`
}

type k8sNodeList struct {
	Items []k8sNodeItem `json:"items"`
}

// ---------------------------------------------------------------------------
// List / Get methods
// ---------------------------------------------------------------------------

// ListNamespaces returns all namespaces.
func (c *K8sClient) ListNamespaces(ctx context.Context) ([]K8sNamespace, error) {
	body, err := c.get(ctx, "/api/v1/namespaces")
	if err != nil {
		return nil, err
	}

	var list k8sNamespaceList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing namespaces: %w", err)
	}

	result := make([]K8sNamespace, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, K8sNamespace{
			Name:   item.Metadata.Name,
			Status: item.Status.Phase,
		})
	}
	return result, nil
}

// ListPods returns pods in a namespace.
func (c *K8sClient) ListPods(ctx context.Context, namespace string) ([]K8sPod, error) {
	body, err := c.get(ctx, fmt.Sprintf("/api/v1/namespaces/%s/pods", namespace))
	if err != nil {
		return nil, err
	}

	var list k8sPodList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing pods: %w", err)
	}

	result := make([]K8sPod, 0, len(list.Items))
	for _, item := range list.Items {
		pod := K8sPod{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Node:      item.Spec.NodeName,
			Status:    podPhaseToStatus(item),
			IP:        item.Status.PodIP,
			Age:       humanAge(item.Metadata.CreationTimestamp),
		}

		// Compute Ready count
		readyCount := 0
		totalCount := len(item.Status.ContainerStatuses)
		totalRestarts := 0
		for _, cs := range item.Status.ContainerStatuses {
			if cs.Ready {
				readyCount++
			}
			totalRestarts += cs.RestartCount

			container := K8sPodContainer{
				Name:         cs.Name,
				Image:        cs.Image,
				Ready:        cs.Ready,
				RestartCount: cs.RestartCount,
				State:        containerStateStr(cs.State.Running, cs.State.Waiting, cs.State.Terminated),
			}
			pod.Containers = append(pod.Containers, container)
		}
		pod.Ready = fmt.Sprintf("%d/%d", readyCount, totalCount)
		pod.Restarts = totalRestarts

		result = append(result, pod)
	}
	return result, nil
}

// ListDeployments returns deployments in a namespace.
func (c *K8sClient) ListDeployments(ctx context.Context, namespace string) ([]K8sDeployment, error) {
	body, err := c.get(ctx, fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments", namespace))
	if err != nil {
		return nil, err
	}

	var list k8sDeploymentList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing deployments: %w", err)
	}

	result := make([]K8sDeployment, 0, len(list.Items))
	for _, item := range list.Items {
		desired := 1
		if item.Spec.Replicas != nil {
			desired = *item.Spec.Replicas
		}

		d := K8sDeployment{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Replicas:  fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, desired),
			Ready:     item.Status.ReadyReplicas,
			UpToDate:  item.Status.UpdatedReplicas,
			Available: item.Status.AvailableReplicas,
			Age:       humanAge(item.Metadata.CreationTimestamp),
		}
		result = append(result, d)
	}
	return result, nil
}

// ListServices returns services in a namespace.
func (c *K8sClient) ListServices(ctx context.Context, namespace string) ([]K8sService, error) {
	body, err := c.get(ctx, fmt.Sprintf("/api/v1/namespaces/%s/services", namespace))
	if err != nil {
		return nil, err
	}

	var list k8sServiceList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing services: %w", err)
	}

	result := make([]K8sService, 0, len(list.Items))
	for _, item := range list.Items {
		svc := K8sService{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Type:      item.Spec.Type,
			ClusterIP: item.Spec.ClusterIP,
			Age:       humanAge(item.Metadata.CreationTimestamp),
		}

		// Ports
		var ports []string
		for _, p := range item.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.Port, p.Protocol))
		}
		svc.Ports = strings.Join(ports, ",")

		// External IP from LoadBalancer ingress
		var extIPs []string
		for _, ing := range item.Status.LoadBalancer.Ingress {
			if ing.IP != "" {
				extIPs = append(extIPs, ing.IP)
			} else if ing.Hostname != "" {
				extIPs = append(extIPs, ing.Hostname)
			}
		}
		if len(extIPs) > 0 {
			svc.ExternalIP = strings.Join(extIPs, ",")
		} else {
			svc.ExternalIP = "<none>"
		}

		result = append(result, svc)
	}
	return result, nil
}

// ListNodes returns all cluster nodes.
func (c *K8sClient) ListNodes(ctx context.Context) ([]K8sNode, error) {
	body, err := c.get(ctx, "/api/v1/nodes")
	if err != nil {
		return nil, err
	}

	var list k8sNodeList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing nodes: %w", err)
	}

	result := make([]K8sNode, 0, len(list.Items))
	for _, item := range list.Items {
		node := K8sNode{
			Name:    item.Metadata.Name,
			Version: item.Status.NodeInfo.KubeletVersion,
			OS:      item.Status.NodeInfo.OperatingSystem,
			Arch:    item.Status.NodeInfo.Architecture,
			Age:     humanAge(item.Metadata.CreationTimestamp),
		}

		// Status from conditions
		node.Status = "NotReady"
		for _, cond := range item.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				node.Status = "Ready"
			}
		}

		// Roles from labels
		var roles []string
		for k := range item.Metadata.Labels {
			if strings.HasPrefix(k, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(k, "node-role.kubernetes.io/")
				if role == "" {
					role = "node"
				}
				roles = append(roles, role)
			}
		}
		if len(roles) == 0 {
			roles = []string{"<none>"}
		}
		node.Roles = strings.Join(roles, ",")

		// Addresses
		for _, addr := range item.Status.Addresses {
			switch addr.Type {
			case "InternalIP":
				node.InternalIP = addr.Address
			case "ExternalIP":
				node.ExternalIP = addr.Address
			}
		}

		// Allocatable resources
		if cpu, ok := item.Status.Allocatable["cpu"]; ok {
			node.CPU = cpu
		}
		if mem, ok := item.Status.Allocatable["memory"]; ok {
			node.Memory = formatK8sMemory(mem)
		}

		result = append(result, node)
	}
	return result, nil
}

// GetPodLogs retrieves the last N lines of logs from a specific container.
func (c *K8sClient) GetPodLogs(ctx context.Context, namespace, podName, containerName string, lines int) (string, error) {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/log?tailLines=%d",
		namespace, podName, lines)
	if containerName != "" {
		path += "&container=" + containerName
	}

	body, err := c.get(ctx, path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func podPhaseToStatus(item k8sPodItem) string {
	// Check for waiting containers first (CrashLoopBackOff, etc.)
	for _, cs := range item.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			return cs.State.Waiting.Reason
		}
		if cs.State.Terminated != nil && cs.State.Terminated.Reason != "" {
			return cs.State.Terminated.Reason
		}
	}
	return item.Status.Phase
}

func containerStateStr(running *struct{}, waiting *struct{ Reason string `json:"reason"` }, terminated *struct{ Reason string `json:"reason"` }) string {
	if running != nil {
		return "Running"
	}
	if waiting != nil {
		if waiting.Reason != "" {
			return waiting.Reason
		}
		return "Waiting"
	}
	if terminated != nil {
		if terminated.Reason != "" {
			return terminated.Reason
		}
		return "Terminated"
	}
	return "Unknown"
}

func humanAge(timestamp string) string {
	if timestamp == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		h := int(d.Hours())
		m := int(math.Mod(d.Minutes(), 60))
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		days := int(d.Hours() / 24)
		if days > 365 {
			return fmt.Sprintf("%dy%dd", days/365, days%365)
		}
		return fmt.Sprintf("%dd", days)
	}
}

// formatK8sMemory converts K8s memory strings like "15885768Ki" to "15Gi".
func formatK8sMemory(mem string) string {
	if strings.HasSuffix(mem, "Ki") {
		val := strings.TrimSuffix(mem, "Ki")
		var ki int
		if _, err := fmt.Sscanf(val, "%d", &ki); err == nil {
			gi := float64(ki) / (1024 * 1024)
			if gi >= 1 {
				return fmt.Sprintf("%.0fGi", gi)
			}
			mi := float64(ki) / 1024
			return fmt.Sprintf("%.0fMi", mi)
		}
	}
	return mem
}
