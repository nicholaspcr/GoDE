package tenant

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	tenant := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "test-tenant-id",
		},
	}

	ctx := NewContext(context.Background(), tenant)

	// Verify tenant was stored in context
	assert.NotNil(t, ctx)

	// Retrieve tenant from context
	retrievedTenant := FromContext(ctx)
	assert.NotNil(t, retrievedTenant)
	assert.Equal(t, "test-tenant-id", retrievedTenant.Ids.TenantId)
}

func TestFromContext_WithTenant(t *testing.T) {
	tenant := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "tenant-123",
		},
	}

	ctx := NewContext(context.Background(), tenant)
	retrievedTenant := FromContext(ctx)

	assert.NotNil(t, retrievedTenant)
	assert.Equal(t, tenant.Ids.TenantId, retrievedTenant.Ids.TenantId)
}

func TestFromContext_WithoutTenant(t *testing.T) {
	// Context without tenant should return default tenant
	ctx := context.Background()
	retrievedTenant := FromContext(ctx)

	assert.NotNil(t, retrievedTenant)
	// Default tenant should be returned
	assert.Equal(t, DefaultTenant, retrievedTenant)
}

func TestFromContext_WithNilTenant(t *testing.T) {
	// Store nil tenant explicitly
	ctx := context.WithValue(context.Background(), tenantKey{}, (*api.Tenant)(nil))
	retrievedTenant := FromContext(ctx)

	// The current implementation returns nil when type assertion succeeds but value is nil
	// This is a known limitation - in production, use NewContext which ensures non-nil tenants
	assert.Nil(t, retrievedTenant)
}

func TestFromContext_WithWrongType(t *testing.T) {
	// Store wrong type in context
	ctx := context.WithValue(context.Background(), tenantKey{}, "not a tenant")
	retrievedTenant := FromContext(ctx)

	// Should return default tenant when type assertion fails
	assert.NotNil(t, retrievedTenant)
	assert.Equal(t, DefaultTenant, retrievedTenant)
}

func TestMultipleTenants(t *testing.T) {
	tenant1 := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "tenant-1",
		},
	}
	tenant2 := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "tenant-2",
		},
	}

	// Create context with first tenant
	ctx1 := NewContext(context.Background(), tenant1)
	retrieved1 := FromContext(ctx1)
	assert.Equal(t, "tenant-1", retrieved1.Ids.TenantId)

	// Create context with second tenant
	ctx2 := NewContext(context.Background(), tenant2)
	retrieved2 := FromContext(ctx2)
	assert.Equal(t, "tenant-2", retrieved2.Ids.TenantId)

	// Verify first context still has first tenant (contexts are immutable)
	retrieved1Again := FromContext(ctx1)
	assert.Equal(t, "tenant-1", retrieved1Again.Ids.TenantId)
}

func TestTenantContextChaining(t *testing.T) {
	tenant1 := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "tenant-1",
		},
	}
	tenant2 := &api.Tenant{
		Ids: &api.TenantIDs{
			TenantId: "tenant-2",
		},
	}

	// Create first context with tenant1
	ctx := NewContext(context.Background(), tenant1)
	assert.Equal(t, "tenant-1", FromContext(ctx).Ids.TenantId)

	// Override with tenant2
	ctx = NewContext(ctx, tenant2)
	assert.Equal(t, "tenant-2", FromContext(ctx).Ids.TenantId)
}

func TestDefaultTenant(t *testing.T) {
	// Verify default tenant exists and has expected properties
	assert.NotNil(t, DefaultTenant)
	assert.NotNil(t, DefaultTenant.Ids)
	assert.NotEmpty(t, DefaultTenant.Ids.TenantId)
}
