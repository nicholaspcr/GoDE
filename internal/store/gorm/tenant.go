package gorm

// TODO: Add tenant operations to inject its value on the context.

// Tenant is a model for the tenant table.
type TenantModel struct {
	Name string `gorm:"primary_key,size:50"`
}
