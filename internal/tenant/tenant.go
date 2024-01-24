// Package tenant provides operations to store tenant information within the context.
package tenant

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

type tenantKey = struct{}

// FromContext returns the Tenant value stored in ctx, if not present it returns an error.
func FromContext(ctx context.Context) *api.Tenant {
	tenant, ok := ctx.Value(tenantKey{}).(*api.Tenant)
	if !ok {
		return DefaultTenant
	}
	return tenant
}

// NewContext returns a new context with the Tenant value stored in it.
func NewContext(ctx context.Context, tenant *api.Tenant) context.Context {
	return context.WithValue(ctx, tenantKey{}, tenant)
}
