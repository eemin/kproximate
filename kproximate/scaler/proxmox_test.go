package scaler

import (
	"regexp"
	"testing"

	"github.com/lupinelab/kproximate/config"
	"github.com/lupinelab/kproximate/kubernetes"
	"github.com/lupinelab/kproximate/proxmox"
	apiv1 "k8s.io/api/core/v1"
)

func TestRequiredScaleEventsFor1CPU(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    1.0,
		Memory: 0,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 0

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 1 {
		t.Errorf("Expected exactly 1 scaleEvent, got: %d", len(requiredScaleEvents))
	}
}

func TestRequiredScaleEventsFor3CPU(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    3.0,
		Memory: 0,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 0

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 2 {
		t.Errorf("Expected exactly 2 scaleEvents, got: %d", len(requiredScaleEvents))
	}
}

func TestRequiredScaleEventsFor1024MBMemory(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    0,
		Memory: 1073741824,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 0

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 1 {
		t.Errorf("Expected exactly 1 scaleEvent, got: %d", len(requiredScaleEvents))
	}
}

func TestRequiredScaleEventsFor3072MBMemory(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    0,
		Memory: 3221225472,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 0

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 2 {
		t.Errorf("Expected exactly 2 scaleEvent, got: %d", len(requiredScaleEvents))
	}
}

func TestRequiredScaleEventsFor1CPU3072MBMemory(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    1,
		Memory: 3221225472,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 0

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 2 {
		t.Errorf("Expected exactly 2 scaleEvent, got: %d", len(requiredScaleEvents))
	}
}

func TestRequiredScaleEventsFor1CPU3072MBMemory1QueuedEvent(t *testing.T) {
	unschedulableResources := kubernetes.UnschedulableResources{
		Cpu:    1,
		Memory: 3221225472,
	}

	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			MaxKpNodes:   3,
		},
	}

	currentEvents := 1

	requiredScaleEvents, err := scaler.requiredScaleEvents(unschedulableResources, currentEvents)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(requiredScaleEvents) != 1 {
		t.Errorf("Expected exactly 1 scaleEvent, got: %d", len(requiredScaleEvents))
	}
}

func TestSelectTargetHosts(t *testing.T) {
	scaler := ProxmoxScaler{
		Proxmox: &proxmox.Mock{
			ClusterStats: []proxmox.HostInformation{
				{
					Id:     "node/host-01",
					Node:   "host-01",
					Cpu:    0.209377325725626,
					Mem:    20394792448,
					Maxmem: 16647962624,
					Status: "online",
				},
				{
					Id:     "node/host-02",
					Node:   "host-02",
					Cpu:    0.209377325725626,
					Mem:    20394792448,
					Maxmem: 16647962624,
					Status: "online",
				},
				{
					Id:     "node/host-03",
					Node:   "host-03",
					Cpu:    0.209377325725626,
					Mem:    11394792448,
					Maxmem: 16647962624,
					Status: "online",
				},
			},
			RunningKpNodes: []proxmox.VmInformation{
				{
					Cpu:     0.114889359119222,
					MaxDisk: 10737418240,
					MaxMem:  2147483648,
					Mem:     1074127542,
					Name:    "kp-node-163c3d58-4c4d-426d-baef-e0c30ecb5fcd",
					NetIn:   35838253204,
					NetOut:  56111331754,
					Node:    "host-03",
					Status:  "running",
					Uptime:  740227,
					VmID:    603,
				},
			},
		},
		config: config.KproximateConfig{
			KpNodeNameRegex:  *regexp.MustCompile(`^kp-node-\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`),
			KpNodeNamePrefix: "kp-node",
		},
	}

	scaleEvents := []*ScaleEvent{
		{
			ScaleType: 1,
			NodeName:  scaler.newKpNodeName(),
		},
		{
			ScaleType: 1,
			NodeName:  scaler.newKpNodeName(),
		},
		{
			ScaleType: 1,
			NodeName:  scaler.newKpNodeName(),
		},
	}

	scaler.SelectTargetHosts(scaleEvents)

	if scaleEvents[0].TargetHost.Node != "host-01" {
		t.Errorf("Expected host-01 to be selected as target host got %s", scaleEvents[0].TargetHost.Node)
	}

	if scaleEvents[1].TargetHost.Node != "host-02" {
		t.Errorf("Expected host-02 to be selected as target host, got %s", scaleEvents[1].TargetHost.Node)
	}

	if scaleEvents[2].TargetHost.Node != "host-03" {
		t.Errorf("Expected host-03 to be selected as target host, got %s", scaleEvents[2].TargetHost.Node)
	}
}

func TestAssessScaleDownForResourceTypeZeroLoad(t *testing.T) {
	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			LoadHeadroom: 0.2,
		},
	}

	scaleDownZeroLoad := scaler.assessScaleDownForResourceType(0, 5, 5)
	if scaleDownZeroLoad == true {
		t.Errorf("Expected false but got %t", scaleDownZeroLoad)
	}
}

func TestAssessScaleDownForResourceTypeAcceptable(t *testing.T) {
	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			LoadHeadroom: 0.2,
		},
	}

	scaleDownAcceptable := scaler.assessScaleDownForResourceType(6, 10, 2)
	if scaleDownAcceptable != true {
		t.Errorf("Expected true but got %t", scaleDownAcceptable)
	}
}

func TestAssessScaleDownForResourceTypeUnAcceptable(t *testing.T) {
	scaler := ProxmoxScaler{
		config: config.KproximateConfig{
			LoadHeadroom: 0.2,
		},
	}

	scaleDownUnAcceptable := scaler.assessScaleDownForResourceType(7, 10, 2)
	if scaleDownUnAcceptable == true {
		t.Errorf("Expected false but got %t", scaleDownUnAcceptable)
	}
}

func TestSelectScaleDownTarget(t *testing.T) {
	node1 := apiv1.Node{}
	node1.Name = "kp-node-163c3d58-4c4d-426d-baef-e0c30ecb5fcd"
	node2 := apiv1.Node{}
	node2.Name = "kp-node-a4f77d63-a944-425d-a980-e7be925b8a6a"
	node3 := apiv1.Node{}
	node3.Name = "kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38"
	kpNodes := []apiv1.Node{
		node1,
		node2,
		node3,
	}

	scaler := ProxmoxScaler{
		Kubernetes: &kubernetes.Mock{
			KpNodes: kpNodes,
			AllocatedResources: map[string]*kubernetes.AllocatedResources{
				"kp-node-163c3d58-4c4d-426d-baef-e0c30ecb5fcd": {
					Cpu:    1.0,
					Memory: 2048.0,
				},
				"kp-node-a4f77d63-a944-425d-a980-e7be925b8a6a": {
					Cpu:    1.0,
					Memory: 2048.0,
				},
				"kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38": {
					Cpu:    1.0,
					Memory: 1048.0,
				},
			},
		},
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 1024,
		},
	}

	scaleEvent := ScaleEvent{
		ScaleType: -1,
	}

	scaler.SelectScaleDownTarget(&scaleEvent)

	if scaleEvent.NodeName != "kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38" {
		t.Errorf("Expected kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38 but got %s", scaleEvent.NodeName)
	}
}

func TestAssessScaleDownIsAcceptable(t *testing.T) {
	scaler := ProxmoxScaler{
		Kubernetes: &kubernetes.Mock{
			AllocatedResources: map[string]*kubernetes.AllocatedResources{
				"kp-node-163c3d58-4c4d-426d-baef-e0c30ecb5fcd": {
					Cpu:    1.0,
					Memory: 1073741824.0,
				},
				"kp-node-a4f77d63-a944-425d-a980-e7be925b8a6a": {
					Cpu:    1.0,
					Memory: 1073741824.0,
				},
				"kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38": {
					Cpu:    1.0,
					Memory: 1073741824.0,
				},
			},
			WorkerNodesAllocatableResources: kubernetes.WorkerNodesAllocatableResources{
				Cpu:    6,
				Memory: 6442450944,
			},
		},
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			LoadHeadroom: 0.2,
		},
	}

	scaleEvent, _ := scaler.AssessScaleDown()

	if scaleEvent.NodeName == "" {
		t.Errorf("scaleEvent had no NodeName")
	}

	if scaleEvent == nil {
		t.Errorf("AssessScaleDown returned nil")
	}

}

func TestAssessScaleDownIsUnacceptable(t *testing.T) {
	scaler := ProxmoxScaler{
		Kubernetes: &kubernetes.Mock{
			AllocatedResources: map[string]*kubernetes.AllocatedResources{
				"kp-node-163c3d58-4c4d-426d-baef-e0c30ecb5fcd": {
					Cpu:    2.0,
					Memory: 2147483648.0,
				},
				"kp-node-a4f77d63-a944-425d-a980-e7be925b8a6a": {
					Cpu:    2.0,
					Memory: 2147483648.0,
				},
				"kp-node-67944692-1de7-4bd0-ac8c-de6dc178cb38": {
					Cpu:    2.0,
					Memory: 2147483648.0,
				},
				"kp-node-a3c5e4ef-4713-473f-b9f7-3abe413c38ff": {
					Cpu:    0.49,
					Memory: 1147483648.0,
				},
				"kp-node-97d74769-22af-420d-9f5e-b2d3c7dd6e7e": {
					Cpu:    1.0,
					Memory: 0.0,
				},
				"kp-node-96f665dd-21c3-4ce1-a1e4-c7717c5338a3": {
					Cpu:    0.0,
					Memory: 0.0,
				},
			},
			WorkerNodesAllocatableResources: kubernetes.WorkerNodesAllocatableResources{
				Cpu:    12,
				Memory: 12884901888,
			},
		},
		config: config.KproximateConfig{
			KpNodeCores:  2,
			KpNodeMemory: 2048,
			LoadHeadroom: 0.2,
		},
	}

	scaleEvent, _ := scaler.AssessScaleDown()

	if scaleEvent != nil {
		t.Errorf("AssessScaleDown did not return nil")
	}
}
