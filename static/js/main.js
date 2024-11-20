document.addEventListener('DOMContentLoaded', function() {
    const deployForm = document.getElementById('deployForm');
    if (deployForm) {
        deployForm.addEventListener('submit', handleDeploy);
    }
});

async function handleDeploy(e) {
    e.preventDefault();

    const formData = new FormData(e.target);
    const deploymentModal = new bootstrap.Modal(document.getElementById('deploymentModal'));
    const statusElement = document.getElementById('deploymentStatus');

    try {
        deploymentModal.show();
        
        const response = await fetch('/deploy', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Erreur lors du déploiement');
        }

        // Polling du statut
        const deploymentId = data.deploymentId;
        pollDeploymentStatus(deploymentId, statusElement);

    } catch (error) {
        statusElement.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
    }
}

async function pollDeploymentStatus(deploymentId, statusElement) {
    try {
        const response = await fetch(`/status/${deploymentId}`);
        const data = await response.json();

        if (data.error) {
            statusElement.innerHTML = `<div class="alert alert-danger">${data.error}</div>`;
            return;
        }

        switch (data.status) {
            case 'terminé':
                statusElement.innerHTML = '<div class="alert alert-success">Déploiement terminé avec succès!</div>';
                break;
            case 'erreur':
                statusElement.innerHTML = `<div class="alert alert-danger">Erreur: ${data.error}</div>`;
                break;
            default:
                statusElement.textContent = 'Déploiement en cours...';
                setTimeout(() => pollDeploymentStatus(deploymentId, statusElement), 5000);
        }
    } catch (error) {
        statusElement.innerHTML = '<div class="alert alert-danger">Erreur lors de la vérification du statut</div>';
    }
}