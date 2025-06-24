package horg

import (
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/datatypes"
)

type OrgM struct {
	hpg.Basic
	Biz       collar.ID      `gorm:"column:biz;index;not null;size:64;"`
	Sovereign bool           `gorm:"column:sovereign;"`
	ID        ID             `gorm:"column:id;primaryKey;"`
	Parent    ID             `gorm:"column:parent;index;"`
	Alias     string         `gorm:"column:alias;size:50"`
	Name      string         `gorm:"column:name;size:200"`
	Tag       []string       `gorm:"column:tag;type:text[];"`
	Meta      datatypes.JSON `gorm:"column:meta;type:json"`
}

func (m *OrgM) TableName() string {
	return "helix_org"
}

func (m *OrgM) ToOrganization() *Organization {
	return &Organization{
		ID:       m.ID,
		Alias:    m.Alias,
		Name:     m.Name,
		Children: make([]*Organization, 0),
	}
}

type OrgMemberM struct {
	hpg.Basic
	Org       ID                    `gorm:"column:org;uniqueIndex:uk_org_member;"`
	Member    sattva.Identification `gorm:"column:member;uniqueIndex:uk_org_member;size:32;"`
	Authority []AuthID              `gorm:"column:authority;serializer:json;"`
}

func (m *OrgMemberM) TableName() string {
	return "helix_org_member"
}

type OrgAuthorityM struct {
	hpg.Basic
	Org     ID       `gorm:"column:org;uniqueIndex:uk_org_id;"`
	InnerID AuthID   `gorm:"column:inner_id;uniqueIndex:uk_org_id;"`
	Name    string   `gorm:"column:name;size:100"`
	Action  []string `gorm:"column:action;serializer:json;"`
}

func (m *OrgAuthorityM) TableName() string {
	return "helix_org_authority"
}
