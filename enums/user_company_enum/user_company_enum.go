package user_company_enum

type Status string

const (
	Active    Status = "active"
	Suspended Status = "suspended"
)

type Role string

const (
	RoleUser           Role = "user"
	RoleCompanyManager Role = "companyManager"
	RoleBackoffice     Role = "backoffice"
)

func (r Role) IsBackoffice() bool {
	return r == RoleBackoffice
}

func (r Role) IsCommonUser() bool {
	return r == RoleUser
}
