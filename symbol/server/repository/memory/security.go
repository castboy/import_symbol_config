package memory

import "import_symbol_config/symbol/server"

type securityRepository struct {
}

var securityRep *securityRepository

func GetSecurityRepository() *securityRepository {
	if securityRep == nil {
		securityRep = &securityRepository{}
	}

	return securityRep
}

// Security operation

// GetSecurityInfo get security info

func (sr *securityRepository) GetSecurityInfo(id int) (*server.Security, error) {
	return nil, nil
}

func (sr *securityRepository) GetAllSecuritiesInfos() ([]*server.Security, error) {
	return nil, nil
}

func (sr *securityRepository) UpdateSymbolSecurity(symbolID, oldSecurityID, newSecurityID int) error {
	return nil
}

func (sr *securityRepository) SetSymbolSecurity(symbolID int, securityID int) error {
	return nil
}

func (sr *securityRepository) SetSecurityInfo(id int, info *server.Security) error {
	return nil
}
