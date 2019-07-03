package server

import (
	"fmt"
	"strings"
	"errors"
)

// Security security structure
type Security struct {
	ID           int      `json:"id" xorm:"id"`
	SecurityName string   `json:"security_name" xorm:"security_name"`
	Description  string   `json:"description" xorm:"description"`
	Symbols      []string `json:"symbols" xorm:"-"`
}

type SecurityOperator interface {
	GetSecurity(id int) (sec *Security, err error)
	GetAllSecurities() (secs []*Security, err error)
	UpdateSecurity(id int, info *Security) error
	InsertSecurity(sec *Security) (err error)
	DeleteSecurity(id int) (err error)
	GetIDByName(securityName string) (ID int, err error)
	GetNameByID(ID int) (securityName string, err error)
	ValidSecurityID(ID int) (valid bool, err error)
	ValidSecurityName(securityName string) (valid bool, err error)
}

type securityOperator struct {
	securityRepo SecurityOperator
}

var securityOp *securityOperator

func GetSecurityOperator() *securityOperator {
	return securityOp
}

func InitSecurityOperator(securityRepo SecurityOperator) *securityOperator {
	if securityOp == nil {
		securityOp = &securityOperator{
			securityRepo,
		}
	}
	return securityOp
}

// GetSecurityInfo get security-structure
func (ss *securityOperator) GetSecurityInfo(id int) (*Security, error) {
	sec, err := ss.securityRepo.GetSecurity(id)
	if err != nil {
		return nil, err
	}

	if sec == nil {
		return nil, errors.New(fmt.Sprintf("invalid security id: %d", id))
	}

	symbols, err := GetSymbolOperator().GetSecuritySymbol(id)
	if err != nil {
		return nil, err
	}

	if symbols == nil {
		// TODO
	}

	sec.Symbols = symbols

	return sec, nil
}

// GetAllSecuritiesInfos get full list of security-structures
func (ss *securityOperator) GetAllSecuritiesInfos() ([]*Security, error) {
	secs, err := ss.securityRepo.GetAllSecurities()
	if err != nil || secs == nil {
		return nil, err
	}

	secSymbolMap, err := GetSymbolOperator().GetAllSecuritySymbols()
	if err != nil {
		return nil, err
	}

	secsInfo := make([]*Security, 0)
	for _, s := range secs {
		sec := new(Security)
		sec.ID = s.ID
		sec.Description = s.Description
		sec.SecurityName = s.SecurityName
		sec.Symbols = strings.Split(secSymbolMap[s.SecurityName], ",")
		secsInfo = append(secsInfo, sec)
	}

	return secsInfo, nil
}

func (so *securityOperator) UpdateSecurityInfo(id int, info *Security) error {
	valid, err := so.securityRepo.ValidSecurityID(id)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New(fmt.Sprintf("invalid security id: %d", id))
	}

	return so.securityRepo.UpdateSecurity(id, info)
}

func (so *securityOperator) DeleteSecurityInfo(id int) error {
	hold, err := GetSymbolOperator().SecurityHoldSymbols(id)
	if err != nil {
		return err
	}
	if hold {
		return errors.New(fmt.Sprintf("can not delete security by id: %d, because it holds symbols", id))
	}

	return so.securityRepo.DeleteSecurity(id)
}

// SetSecurityInfo Only available to set securityName and description, but not for symbols list
func (so *securityOperator) InsertSecurityInfo(info *Security) error {
	return so.securityRepo.InsertSecurity(info)
}

func (ss *securityOperator) GetIDByName(securityName string) (ID int, err error) {
	return ss.securityRepo.GetIDByName(securityName)
}

func (ss *securityOperator) GetNameByID(ID int) (securityName string, err error) {
	return ss.securityRepo.GetNameByID(ID)
}

func (ss *securityOperator) ValidSecurityID(ID int) (valid bool, err error) {
	return ss.securityRepo.ValidSecurityID(ID)
}

func (ss *securityOperator) ValidSecurityName(securityName string) (valid bool, err error) {
	return ss.securityRepo.ValidSecurityName(securityName)
}
