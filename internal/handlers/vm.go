package handlers

import (
	"azure-vm-deployer/internal/models"
	"azure-vm-deployer/internal/terraform"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeployPage(c *gin.Context) {
	c.HTML(http.StatusOK, "deploy.html", gin.H{
		"title": "Déployer une VM",
	})
}

func DeployVM(c *gin.Context) {
	var req models.VMRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation des entrées
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Générer la configuration Terraform
	tfConfig, err := terraform.GenerateConfig(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de génération de configuration"})
		return
	}

	// Lancer le déploiement de manière asynchrone
	deploymentID := terraform.StartDeployment(tfConfig)

	c.JSON(http.StatusOK, gin.H{
		"message": "Déploiement en cours",
		"deploymentId": deploymentID,
	})
}

func VMStatus(c *gin.Context) {
	deploymentID := c.Param("id")
	status := terraform.GetDeploymentStatus(deploymentID)
	c.JSON(http.StatusOK, status)
}

func ListVMs(c *gin.Context) {
	vms := terraform.ListDeployedVMs()
	c.JSON(http.StatusOK, vms)
}