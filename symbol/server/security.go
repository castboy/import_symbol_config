package server

import (
	"github.com/juju/errors"
	"strconv"
	"strings"
)

// Security security structure
type Security struct {
	ID           int      `json:"id" xorm:"id"`
	SecurityName string   `json:"security_name" xorm:"security_name"`
	Description  string   `json:"description" xorm:"description"`
	Symbols      []string `json:"symbols" xorm:"-"`
}

type SecurityOperator interface {
	GetSecurity(id int) (sec *Security, exist bool, err error)
	GetAllSecurities() (secs []*Security, err error)
	UpdateSecurity(id int, info *Security) error
	InsertSecurity(sec *Security) (err error)
	DeleteSecurity(id int) (err error)
	GetIDByName(securityName string) (ID int, exist bool, err error)
	GetNameByID(ID int) (securityName string, exist bool, err error)
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
	sec, exist, err := ss.securityRepo.GetSecurity(id)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return nil, err
	}

	if !exist {
		err = errors.NotFoundf("security id %d", id)
		return nil, err
	}

	symbols, err := GetSymbolOperator().GetSecuritySymbol(id)
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return nil, err
	}

	if len(symbols) != 0 {
		sec.Symbols = symbols
	}

	return sec, nil
}

// GetAllSecuritiesInfos get full list of security-structures
func (ss *securityOperator) GetAllSecuritiesInfos() ([]*Security, error) {
	secs, err := ss.securityRepo.GetAllSecurities()
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return nil, err
	}

	if len(secs) == 0 {
		err = errors.NotFoundf("securities")
		return nil, err
	}

	secSymbs, err := GetSymbolOperator().GetAllSecuritySymbols()
	if err != nil {
		err = errors.Annotatef(err, "sql exec")
		return nil, err
	}

	if len(secSymbs) == 0 {
		err = errors.NotFoundf("securities")
		return nil, err
	}

	secSymbols := make(map[int]string)
	for i := range secSymbs {
		securityID, err := strconv.Atoi(secSymbs[i]["security_id"])
		if err != nil {
			// TODO
			continue
		}

		secSymbols[securityID] = secSymbs[i]["symbol"]
	}

	secsInfo := make([]*Security, 0)
	for _, s := range secs {
		sec := new(Security)
		sec.ID = s.ID
		sec.Description = s.Description
		sec.SecurityName = s.SecurityName

		if _, ok := secSymbols[s.ID]; ok {
			sec.Symbols = strings.Split(secSymbols[s.ID], ",")
		}

		secsInfo = append(secsInfo, sec)
	}

	return secsInfo, nil
}

func (so *securityOperator) UpdateSecurityInfo(id int, info *Security) error {
	valid, err := so.securityRepo.ValidSecurityID(id)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	if !valid {
		return errors.NotValidf("security id %d", id)
	}

	err = so.securityRepo.UpdateSecurity(id, info)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	return nil
}

func (so *securityOperator) DeleteSecurityInfo(id int) error {
	hold, err := GetSymbolOperator().SecurityHoldSymbols(id)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	if hold {
		return errors.Forbiddenf("security id %d", id)
	}

	err = so.securityRepo.DeleteSecurity(id)
	if err != nil {
		return errors.Annotatef(err, "sql exec")
	}

	return nil
}

// SetSecurityInfo Only available to set securityName and description, but not for symbols list
func (so *securityOperator) InsertSecurityInfo(info *Security) error {
	return so.securityRepo.InsertSecurity(info)
}

func (ss *securityOperator) GetIDByName(securityName string) (ID int, exist bool, err error) {
	return ss.securityRepo.GetIDByName(securityName)
}

func (ss *securityOperator) GetNameByID(ID int) (securityName string, exist bool, err error) {
	return ss.securityRepo.GetNameByID(ID)
}

func (ss *securityOperator) ValidSecurityID(ID int) (valid bool, err error) {
	return ss.securityRepo.ValidSecurityID(ID)
}

func (ss *securityOperator) ValidSecurityName(securityName string) (valid bool, err error) {
	return ss.securityRepo.ValidSecurityName(securityName)
}
