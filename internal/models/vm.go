package models

import "errors"

type VMRequest struct {
	Name     string `json:"vmName" form:"vmName"`
	Size     string `json:"vmSize" form:"vmSize"`
	Location string `json:"location" form:"location"`
	OS       string `json:"os" form:"os"`
}

func (v *VMRequest) Validate() error {
	if v.Name == "" {
		return errors.New("le nom de la VM est requis")
	}

	validSizes := map[string]bool{
		"Standard_B1s":  true,
		"Standard_B2s":  true,
		"Standard_B4ms": true,
	}
	if !validSizes[v.Size] {
		return errors.New("taille de VM invalide")
	}

	validLocations := map[string]bool{
		"westeurope":    true,
		"francecentral": true,
		"northeurope":   true,
	}
	if !validLocations[v.Location] {
		return errors.New("région invalide")
	}

	validOS := map[string]bool{
		"windows2019": true,
		"ubuntu2004":  true,
	}
	if !validOS[v.OS] {
		return errors.New("système d'exploitation invalide")
	}

	return nil
}