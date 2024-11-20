package terraform

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	deployments = make(map[string]*Deployment)
	mu          sync.RWMutex
)

type Deployment struct {
	ID        string
	Status    string
	StartTime time.Time
	Error     string
}

func StartDeployment(config string) string {
	deploymentID := fmt.Sprintf("deploy-%d", time.Now().Unix())
	
	deployment := &Deployment{
		ID:        deploymentID,
		Status:    "en cours",
		StartTime: time.Now(),
	}

	mu.Lock()
	deployments[deploymentID] = deployment
	mu.Unlock()

	go func() {
		workingDir := filepath.Join(os.TempDir(), deploymentID)
		if err := os.MkdirAll(workingDir, 0755); err != nil {
			updateDeploymentStatus(deploymentID, "erreur", err.Error())
			return
		}

		if err := os.WriteFile(filepath.Join(workingDir, "main.tf"), []byte(config), 0644); err != nil {
			updateDeploymentStatus(deploymentID, "erreur", err.Error())
			return
		}

		tf, err := tfexec.NewTerraform(workingDir, "terraform")
		if err != nil {
			updateDeploymentStatus(deploymentID, "erreur", err.Error())
			return
		}

		if err := tf.Init(context.Background()); err != nil {
			updateDeploymentStatus(deploymentID, "erreur", err.Error())
			return
		}

		if err := tf.Apply(context.Background()); err != nil {
			updateDeploymentStatus(deploymentID, "erreur", err.Error())
			return
		}

		updateDeploymentStatus(deploymentID, "terminé", "")
	}()

	return deploymentID
}

func updateDeploymentStatus(id, status, errorMsg string) {
	mu.Lock()
	defer mu.Unlock()

	if deployment, exists := deployments[id]; exists {
		deployment.Status = status
		deployment.Error = errorMsg
	}
}

func GetDeploymentStatus(id string) map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()

	if deployment, exists := deployments[id]; exists {
		return map[string]interface{}{
			"id":         deployment.ID,
			"status":     deployment.Status,
			"startTime":  deployment.StartTime,
			"error":      deployment.Error,
		}
	}

	return map[string]interface{}{
		"error": "Déploiement non trouvé",
	}
}

func ListDeployedVMs() []map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()

	var vms []map[string]interface{}
	for _, deployment := range deployments {
		vms = append(vms, map[string]interface{}{
			"id":        deployment.ID,
			"status":    deployment.Status,
			"startTime": deployment.StartTime,
		})
	}
	return vms
}