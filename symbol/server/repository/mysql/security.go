package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"import_symbol_config/symbol/server"
)

type securityRepository struct {
	engine *xorm.Engine
}

var securityRep *securityRepository

func GetSecurityRepository() *securityRepository {
	if securityRep == nil {
		securityRep = &securityRepository{
			engine: xEngine,
		}
	}

	return securityRep
}

func (sr *securityRepository) GetSecurity(id int) (sec *server.Security, exist bool, err error) {
	sec = new(server.Security)
	exist, err = sr.engine.Table(server.Security{}).Where("id=?", id).Get(sec)
	return
}

// GetAllSecuritiesInfo get all securitites info
func (sr *securityRepository) GetAllSecurities() (secs []*server.Security, err error) {
	err = sr.engine.Table(server.Security{}).OrderBy("id").Find(&secs)
	return

}

// InsertSecurity insert security info into db
func (sr *securityRepository) InsertSecurity(sec *server.Security) (err error) {
	_, err = sr.engine.Table(server.Security{}).Omit("id").Insert(sec)
	return
}

// UpdateSecurity update security info
func (sr *securityRepository) UpdateSecurity(id int, info *server.Security) (err error) {
	_, err = sr.engine.Table(server.Security{}).Where("id=?", id).AllCols().Update(info)
	return
}

// DeleteSecurity delete security info
func (sr *securityRepository) DeleteSecurity(id int) (err error) {
	_, err = sr.engine.Table(server.Security{}).Where("id=?", id).Delete(server.Security{})
	return
}

func (sr *securityRepository) GetIDByName(securityName string) (ID int, exist bool, err error) {
	exist, err = sr.engine.Table(server.Security{}).Select("id").Where("security_name=?", securityName).Get(&ID)
	return
}

func (sr *securityRepository) GetNameByID(ID int) (securityName string, exist bool, err error) {
	exist, err = sr.engine.Table(server.Security{}).Select("security_name").Where("id=?", ID).Get(&securityName)
	return
}

func (sr *securityRepository) ValidSecurityID(ID int) (valid bool, err error) {
	return sr.engine.Table(server.Security{}).Where("id=?", ID).Exist()
}

func (sr *securityRepository) ValidSecurityName(securityName string) (valid bool, err error) {
	return sr.engine.Table(server.Security{}).Where("security_name=?", securityName).Exist()
}
