package projectcfg

type MemberSlug string

type Members []*Member

type MemberType string

type Member struct {
	Name     string     `json:"name"`
	Slug     string     `json:"slug,omitempty"`
	Type     MemberType `json:"type"`
	Roles    RoleTypes  `json:"roles"`
	Email    string     `json:"email"`
	Website  string     `json:"website,omitempty"`
	LinkedIn string     `json:"linked_in,omitempty"`
	Twitter  string     `json:"twitter,omitempty"`
	Blog     string     `json:"blog,omitempty"`
	Details    Details      `json:"name,omitempty"`
}

const (
	PersonMember       MemberType = "person"
	OrganizationMember MemberType = "org"
	ProjectMember      MemberType = "project"
	OtherMember        MemberType = "other"
)

type RoleTypes []RoleType
type RoleType string

const (
	UserRole      RoleType = "user"
	TestRole      RoleType = "test"
	DocRole       RoleType = "doc"
	DevRole       RoleType = "dev"
	DesignRole    RoleType = "design"
	UxRole        RoleType = "ux"
	ContentRole   RoleType = "content"
	DevOpsRole    RoleType = "devops"
	SysAdminRole  RoleType = "sysadmin"
	ProductRole   RoleType = "product"
	PmRole        RoleType = "pm"
	ArchitectRole RoleType = "architect"
	AdvocateRole  RoleType = "advocate"
	SponsorRole   RoleType = "sponsor"
	OwnerRole     RoleType = "owner"
	OtherRole     RoleType = "other"
)
