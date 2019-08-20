package main

type VaultService struct {
	label       string
	name        string
	credentials Credentials
}

type Credentials struct {
	address         string
	auth            Auth
	backends        Backends
	backends_shared Backends_Shared
}

type Auth struct {
	token string
}

type Backends struct {
	generic string
	transit string
}

type Backends_Shared struct {
	organization string
	space        string
}

func (v VaultService) getTokenfromCF(data interface{}) VaultService {
	m := data.(map[string]interface{})
	if token, ok := m["token"].(string); ok {
		v.credentials.auth.token = token
	}
	return v
}

func (v VaultService) getBackendsfromCF(data interface{}) VaultService {
	m := data.(map[string]interface{})
	if generic, ok := m["generic"].(string); ok {
		v.credentials.backends.generic = generic
	}
	if transit, ok := m["transit"].(string); ok {
		v.credentials.backends.transit = transit
	}
	return v
}

func (v VaultService) getBackendsSharedfromCF(data interface{}) VaultService {
	m := data.(map[string]interface{})
	if organization, ok := m["organization"].(string); ok {
		v.credentials.backends_shared.organization = organization
	}
	if space, ok := m["space"].(string); ok {
		v.credentials.backends_shared.space = space
	}
	return v
}
