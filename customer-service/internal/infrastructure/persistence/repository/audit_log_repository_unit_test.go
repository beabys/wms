package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditLogRepositoryTypeCheck(t *testing.T) {
	var r *AuditLogRepository
	assert.Nil(t, r)
}
