/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package model

type CookieTTL struct {
	AuthTTL int `json:"authTTL"`
}

type WhiteList struct {
	AllowList []string `json:"allowList"`
}

type Department struct {
	Aliases      []string `json:"aliases"`
	CreatedAt    string   `json:"createdAt"`
	Description  string   `json:"description"`
	Email        string   `json:"email"`
	ExternalId   string   `json:"externalId"`
	HeadId       string   `json:"headId"`
	Id           int      `json:"id"`
	Label        string   `json:"label"`
	MembersCount int      `json:"membersCount"`
	Name         string   `json:"name"`
	ParentId     int      `json:"parentId"`
	Removed      bool     `json:"removed"`
}

type DepartmentList struct {
	Departments []Department `json:"departments"`
	Page        int          `json:"page"`
	Pages       int          `json:"pages"`
	PerPage     int          `json:"perPage"`
	Total       int          `json:"total"`
}

type AuditEvent struct {
	ClientIp             string `json:"clientIP"`
	Date                 string `json:"date"`
	EventType            string `json:"eventType"`
	LastModificationDate string `json:"lastModificationDate"`
	OrgId                int    `json:"orgId"`
	OwnerLogin           string `json:"ownerId"`
	OwnerName            string `json:"ownerName"`
	OwnerUid             string `json:"ownerUid"`
	Path                 string `json:"path"`
	RequestId            string `json:"requestId"`
	ResourceFileId       string `json:"resourceFileId"`
	Rights               string `json:"rights"`
	Size                 string `json:"size"`
	UniqId               string `json:"uniqId"`
	UserLogin            string `json:"userLogin"`
	UserName             string `json:"userName"`
	UserUid              string `json:"userUid"`
}

type DiskAuditLog struct {
	Events        []AuditEvent `json:"events"`
	NextPageToken string       `json:"nextPageToken"`
}

type MfaSetup struct {
	Duration  int    `json:"duration"`
	Enabled   bool   `json:"enabled"`
	EnabledAt string `json:"enabledAt"`
}

type MfaActivation struct {
	Duration         int    `json:"duration"`
	LogoutUsers      bool   `json:"logoutUsers"`
	ValidationMethod string `json:"validationMethod"`
}

type Organization struct {
	Email            string `json:"email"`
	Fax              string `json:"fax"`
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	Language         string `json:"language"`
	SubscriptionPlan string `json:"subscriptionPlan"`
}

type OrganizationList struct {
	NextPage      string         `json:"nextPageToken"`
	Organizations []Organization `json:"organizations"`
}

type ContactInfo struct {
	Alias     bool   `json:"alias"`
	Label     string `json:"label"`
	Main      bool   `json:"main"`
	Synthetic bool   `json:"synthetic"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type UserName struct {
	First  string `json:"first"`
	Last   string `json:"last"`
	Middle string `json:"middle"`
}

type User struct {
	About        string        `json:"about"`
	Aliases      []string      `json:"aliases"`
	AvatarId     string        `json:"avatarId"`
	Birthday     string        `json:"birthday"`
	Contacts     []ContactInfo `json:"contacts"`
	CreatedAt    string        `json:"createdAt"`
	DepartmentId int           `json:"departmentId"`
	DisplayName  string        `json:"displayName"`
	Email        string        `json:"email"`
	ExternalId   string        `json:"externalId"`
	Gender       string        `json:"gender"`
	Groups       []int         `json:"groups"`
	Id           string        `json:"id"`
	IsAdmin      bool          `json:"isAdmin"`
	IsDismissed  bool          `json:"isDismissed"`
	IsEnabled    bool          `json:"isEnabled"`
	IsRobot      bool          `json:"isRobot"`
	Language     string        `json:"language"`
	Name         UserName      `json:"name"`
	Nickname     string        `json:"nickname"`
	Position     string        `json:"position"`
	Timezone     string        `json:"timezone"`
	UpdatedAt    string        `json:"updatedAt"`
}

type UserList struct {
	Page    int    `json:"page"`
	Pages   int    `json:"pages"`
	PerPage int    `json:"perPage"`
	Total   int    `json:"total"`
	Users   []User `json:"users"`
}

type RmAliasResponse struct {
	Alias   string `json:"alias"`
	Removed bool   `json:"removed"`
}

type UserMfaSetup struct {
	UserId           string `json:"userId"`
	Has2fa           bool   `json:"has2fa"`
	HasSecurityPhone bool   `json:"hasSecurityPhone"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

type Group struct {
	AdminIds     []string `json:"adminIds,omitempty"`
	Aliases      []string `json:"aliases,omitempty"`
	AuthorId     string   `json:"authorId,omitempty"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	Description  string   `json:"description,omitempty"`
	Email        string   `json:"email,omitempty"`
	ExternalId   string   `json:"externalId,omitempty"`
	Id           int      `json:"id"`
	Label        string   `json:"label,omitempty"`
	MemberOf     []int    `json:"memberOf,omitempty"`
	Members      []Member `json:"members,omitempty"`
	MembersCount int      `json:"membersCount,omitempty"`
	Name         string   `json:"name"`
	Removed      bool     `json:"removed,omitempty"`
	Type         string   `json:"type,omitempty"`
}

type Member struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type GroupList struct {
	Groups  []Group `json:"groups"`
	Page    int     `json:"page"`
	Pages   int     `json:"pages"`
	PerPage int     `json:"perPage"`
	Total   int     `json:"total"`
}

type UserSign struct {
	Emails    []string `json:"emails"`
	IsDefault bool     `json:"isDefault"`
	Lang      string   `json:"lang"`
	Text      string   `json:"text"`
}

type UserSenderInfo struct {
	DefaultFrom  string     `json:"defaultFrom"`
	FromName     string     `json:"fromName"`
	SignPosition string     `json:"signPosition"`
	Signs        []UserSign `json:"signs"`
}

type AutoreplyRule struct {
	RuleId   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	Text     string `json:"text"`
}

type ForwardRule struct {
	Address   string `json:"address"`
	RuleId    int    `json:"ruleId"`
	RuleName  string `json:"ruleName"`
	WithStore bool   `json:"withStore"`
}

type UserMailRules struct {
	Autoreplies []AutoreplyRule `json:"autoreplies"`
	Forwards    []ForwardRule   `json:"forwards"`
}

type GroupMember struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	MembersCount int    `json:"membersCount"`
}

type GroupMemberList struct {
	Departments []GroupMember `json:"departments"`
	Groups      []GroupMember `json:"groups"`
	Users       []User        `json:"users"`
}

type GroupMemberResponse struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Deleted bool   `json:"deleted"` // if removed
	Added   bool   `json:"added"`   // if added
}

type DnsRecord struct {
	Address    string `json:"address"`
	Exchange   string `json:"exchange"`
	Flag       int    `json:"flag"`
	Name       string `json:"name"`
	Port       int    `json:"port"`
	Preference int    `json:"preference"`
	Priority   int    `json:"priority"`
	RecordId   int    `json:"recordId"`
	Tag        string `json:"tag"`
	Target     string `json:"target"`
	Text       string `json:"text"`
	Ttl        int    `json:"ttl"`
	Type       string `json:"type"`
	Value      string `json:"value"`
	Weight     int    `json:"weight"`
}

type DnsRecordList struct {
	Page    int         `json:"page"`
	Pages   int         `json:"pages"`
	PerPage int         `json:"perPage"`
	Records []DnsRecord `json:"records"`
	Total   int         `json:"total"`
}

type StatusCheck struct {
	Match bool   `json:"match"`
	Value string `json:"value"`
}

type DomainStatus struct {
	DKIM      StatusCheck `json:"dkim"`
	LastAdded string      `json:"lastAdded"`
	LastCheck string      `json:"lastCheck"`
	MX        StatusCheck `json:"mx"`
	Name      string      `json:"name"`
	NS        StatusCheck `json:"ns"`
	SPF       StatusCheck `json:"spf"`
}

type Domain struct {
	Country   string       `json:"country"`
	Delegated bool         `json:"delegated"`
	Master    bool         `json:"master"`
	MX        bool         `json:"mx"`
	Name      string       `json:"name"`
	Status    DomainStatus `json:"status"`
	Verified  bool         `json:"verified"`
}

type DomainList struct {
	Domains []Domain `json:"domains"`
	Page    int      `json:"page"`
	Pages   int      `json:"pages"`
	PerPage int      `json:"perPage"`
	Total   int      `json:"total"`
}

type AdminList struct {
	AdminIds []string `json:"adminIds"`
}

type ContactInfoList struct {
	Items []ContactInfo `json:"contacts"`
}

type MailAccessSettings struct {
	Items []string `json:"roles"`
}

type MailAccessResponse struct {
	TaskId string `json:"taskId"`
}

type TaskStatusResponse struct {
	Status string `json:"status"`
}

type Resource struct {
	ResourceId   string   `json:"resourceId"`
	ResourceType string   `json:"type"`
	Items        []string `json:"roles"`
}

type ResourceList struct {
	Items []Resource `json:"resources"`
}

type Actor struct {
	ActorId string   `json:"actorId"`
	Items   []string `json:"roles"`
}

type ActorList struct {
	Items []Actor `json:"actors"`
}

type MailboxListItem struct {
	ResourceId string `json:"resourceId"`
	Count      int    `json:"count"`
}

type MailboxList struct {
	Items   []MailboxListItem `json:"resources"`
	Page    int               `json:"page"`
	PerPage int               `json:"perPage"`
	Total   int               `json:"total"`
}

type MailboxInfo struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type UserDeletionResponse struct {
	Deleted bool   `json:"deleted"`
	UserId  string `json:"userId"`
}

type OAuthStatusResponse struct {
	Restricted bool `json:"restricted"`
}
