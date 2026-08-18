package rbac

type RoleInput struct {
	Name        string `json:"name" binding:"required,min=1,max=80"`
	Code        string `json:"code" binding:"required,min=2,max=80"`
	Description string `json:"description" binding:"max=500"`
	Enabled     *bool  `json:"enabled"`
}

type PermissionInput struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Code        string `json:"code" binding:"required,min=3,max=120"`
	Module      string `json:"module" binding:"required,min=1,max=60"`
	Description string `json:"description" binding:"max=500"`
	Enabled     *bool  `json:"enabled"`
	GroupID     *int   `json:"group_id" binding:"omitempty,min=1"`
}

type MenuInput struct {
	Name      string `json:"name" binding:"required,min=1,max=100"`
	Path      string `json:"path" binding:"max=240"`
	Icon      string `json:"icon" binding:"max=100"`
	Component string `json:"component" binding:"max=200"`
	Redirect  string `json:"redirect" binding:"max=240"`
	Type      string `json:"type" binding:"required,oneof=CATALOG MENU BUTTON EMBEDDED LINK"`
	ParentID  *int   `json:"parent_id" binding:"omitempty,min=1"`
	SortOrder int    `json:"sort_order"`
	Enabled   *bool  `json:"enabled"`
}

type GrantInput struct {
	PermissionIDs []int `json:"permission_ids" binding:"dive,min=1"`
	MenuIDs       []int `json:"menu_ids" binding:"dive,min=1"`
}

type PermissionGroupInput struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Module      string `json:"module" binding:"required,min=1,max=60"`
	Description string `json:"description" binding:"max=500"`
	ParentID    *int   `json:"parent_id" binding:"omitempty,min=1"`
	SortOrder   int    `json:"sort_order"`
	Enabled     *bool  `json:"enabled"`
}
type APIInput struct {
	Name        string `json:"name" binding:"required,min=1,max=160"`
	Method      string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	Path        string `json:"path" binding:"required,max=240"`
	Description string `json:"description" binding:"max=500"`
	Enabled     *bool  `json:"enabled"`
}
type ResourceGrantInput struct {
	APIIDs  []int `json:"api_ids" binding:"dive,min=1"`
	MenuIDs []int `json:"menu_ids" binding:"dive,min=1"`
}

type RoleView struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	Enabled       bool   `json:"enabled"`
	IsSystem      bool   `json:"is_system"`
	PermissionIDs []int  `json:"permission_ids"`
	MenuIDs       []int  `json:"menu_ids"`
}

type PermissionView struct {
	ID          int    `json:"id"`
	RoleCount   int    `json:"role_count"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Module      string `json:"module"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	GroupID     *int   `json:"group_id"`
	APIIDs      []int  `json:"api_ids"`
	MenuIDs     []int  `json:"menu_ids"`
}

type MenuView struct {
	ID        int    `json:"id"`
	SortOrder int    `json:"sort_order"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Component string `json:"component"`
	Redirect  string `json:"redirect"`
	Type      string `json:"type"`
	ParentID  *int   `json:"parent_id"`
	Enabled   bool   `json:"enabled"`
}

type PermissionGroupView struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	Enabled     bool   `json:"enabled"`
}
type APIView struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
type AuditView struct {
	ID         int    `json:"id"`
	OperatorID int    `json:"operator_id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   int    `json:"target_id"`
	RequestID  string `json:"request_id"`
	IPAddress  string `json:"ip_address"`
	CreateTime string `json:"create_time"`
}

type Snapshot struct {
	Roles       []RoleView            `json:"roles"`
	Permissions []PermissionView      `json:"permissions"`
	Menus       []MenuView            `json:"menus"`
	Groups      []PermissionGroupView `json:"groups"`
	APIs        []APIView             `json:"apis"`
	AuditLogs   []AuditView           `json:"audit_logs"`
}

type AccessView struct {
	UserID      int        `json:"user_id"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions"`
	MenuIDs     []int      `json:"menu_ids"`
	Menus       []MenuView `json:"menus"`
}
