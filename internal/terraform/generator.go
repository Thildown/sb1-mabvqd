package terraform

import (
	"azure-vm-deployer/internal/models"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"text/template"
)

func generatePassword() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GenerateConfig(req models.VMRequest) (string, error) {
	const tfTemplate = `
provider "azurerm" {
    features {}
}

resource "azurerm_resource_group" "rg" {
    name     = "rg-{{.Name}}"
    location = "{{.Location}}"
    tags = {
        environment = "production"
        deployed_by = "vm-deployer"
    }
}

resource "azurerm_virtual_network" "vnet" {
    name                = "vnet-{{.Name}}"
    address_space       = ["10.0.0.0/16"]
    location           = azurerm_resource_group.rg.location
    resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_subnet" "subnet" {
    name                 = "subnet-{{.Name}}"
    resource_group_name  = azurerm_resource_group.rg.name
    virtual_network_name = azurerm_virtual_network.vnet.name
    address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_interface" "nic" {
    name                = "nic-{{.Name}}"
    location            = azurerm_resource_group.rg.location
    resource_group_name = azurerm_resource_group.rg.name

    ip_configuration {
        name                          = "internal"
        subnet_id                     = azurerm_subnet.subnet.id
        private_ip_address_allocation = "Dynamic"
    }
}

resource "azurerm_virtual_machine" "vm" {
    name                  = "{{.Name}}"
    location             = azurerm_resource_group.rg.location
    resource_group_name  = azurerm_resource_group.rg.name
    network_interface_ids = [azurerm_network_interface.nic.id]
    vm_size              = "{{.Size}}"

    storage_image_reference {
        {{if eq .OS "windows2019"}}
        publisher = "MicrosoftWindowsServer"
        offer     = "WindowsServer"
        sku       = "2019-Datacenter"
        {{else}}
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "20.04-LTS"
        {{end}}
        version   = "latest"
    }

    storage_os_disk {
        name              = "osdisk-{{.Name}}"
        caching           = "ReadWrite"
        create_option     = "FromImage"
        managed_disk_type = "Standard_LRS"
    }

    os_profile {
        computer_name  = "{{.Name}}"
        admin_username = "vmadmin"
        admin_password = "{{.Password}}"
    }

    {{if eq .OS "windows2019"}}
    os_profile_windows_config {
        provision_vm_agent = true
    }
    {{else}}
    os_profile_linux_config {
        disable_password_authentication = false
    }
    {{end}}
}
`

	tmpl, err := template.New("terraform").Parse(tfTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		models.VMRequest
		Password string
	}{
		VMRequest: req,
		Password:  generatePassword(),
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}