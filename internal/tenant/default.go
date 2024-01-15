package tenant

import "github.com/nicholaspcr/GoDE/pkg/api/v1"

// DefaultTenant is the default tenant value.
var DefaultTenant = &api.Tenant{Ids: &api.TenantIDs{TenantId: "default"}}
