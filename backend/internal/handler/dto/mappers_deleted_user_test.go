//go:build unit || !integration

package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/LuckyKuang/sub2api-plus/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUserFromServiceShallow_MapsDeletedAt(t *testing.T) {
	ts := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)

	deleted := UserFromServiceShallow(&service.User{ID: 1, Email: "d@test.com", DeletedAt: &ts})
	require.NotNil(t, deleted.DeletedAt)
	require.Equal(t, ts, *deleted.DeletedAt)

	active := UserFromServiceShallow(&service.User{ID: 2, Email: "a@test.com"})
	require.Nil(t, active.DeletedAt, "active user must have nil DeletedAt")
}

func TestAdminUserReputationScoreZeroIsSerialized(t *testing.T) {
	out := UserFromServiceAdmin(&service.User{ID: 9, Email: "zero@test.com", ReputationScore: 0})
	require.NotNil(t, out)
	AttachReputationTier(&out.User, service.ReputationSettings{
		Enabled: true,
		Tiers: []service.ReputationTier{
			{ID: "tier-pass", Name: "pass", Score: 60},
			{ID: "tier-risk", Name: "risk", Score: 40},
		},
	})
	raw, err := json.Marshal(out)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"reputation_score":0`)
	require.Contains(t, string(raw), `"id":"tier-risk"`)
}

func TestUserFacingProfileOmitsReputationScore(t *testing.T) {
	out := UserFromServiceShallow(&service.User{ID: 9, Email: "user@test.com", ReputationScore: 60})
	require.NotNil(t, out)
	require.Nil(t, out.ReputationScore)
	raw, err := json.Marshal(out)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "reputation_score")
	require.NotContains(t, string(raw), "reputation_tier")
}

func TestAttachReputationTierDisabledOmitsScore(t *testing.T) {
	out := UserFromServiceAdmin(&service.User{ID: 9, Email: "zero@test.com", ReputationScore: 60})
	require.NotNil(t, out)
	AttachReputationTier(&out.User, service.ReputationSettings{Enabled: false})
	raw, err := json.Marshal(out)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "reputation_score")
	require.NotContains(t, string(raw), "reputation_tier")
}
