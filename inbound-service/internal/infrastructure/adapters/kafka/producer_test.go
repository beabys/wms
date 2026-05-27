package kafka

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/inbound-service/internal/domain/inbound/model"
)

func TestEventPublisher_EventPayloads(t *testing.T) {
	t.Run("submitted event serialization", func(t *testing.T) {
		in := &model.Inbound{
			ID:           "inb_1",
			CustomerID:   "cust_1",
			Items:        []model.InboundItem{{SKU: "SKU001", QuantityDeclared: 10}},
			ExpectedDate: "2026-06-01",
		}

		evt := SubmittedEvent{
			InboundID:    in.ID,
			CustomerID:   in.CustomerID,
			ItemsCount:   len(in.Items),
			ExpectedDate: in.ExpectedDate,
			Timestamp:    time.Now(),
		}

		data, err := json.Marshal(evt)
		require.NoError(t, err)

		var decoded SubmittedEvent
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "inb_1", decoded.InboundID)
		assert.Equal(t, "cust_1", decoded.CustomerID)
		assert.Equal(t, 1, decoded.ItemsCount)
	})

	t.Run("approved event serialization", func(t *testing.T) {
		in := &model.Inbound{
			ID:         "inb_2",
			CustomerID: "cust_1",
			Items: []model.InboundItem{
				{SKU: "SKU001", QuantityDeclared: 10, QuantityReceived: 10, Weight: 5.0},
			},
		}

		items := make([]ApprovedItem, len(in.Items))
		for i, item := range in.Items {
			items[i] = ApprovedItem{
				SKU:              item.SKU,
				QuantityDeclared: item.QuantityDeclared,
				QuantityReceived: item.QuantityReceived,
				Weight:           item.Weight,
			}
		}

		evt := ApprovedEvent{
			InboundID:  in.ID,
			CustomerID: in.CustomerID,
			Items:      items,
			Timestamp:  time.Now(),
		}

		data, err := json.Marshal(evt)
		require.NoError(t, err)

		var decoded ApprovedEvent
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "inb_2", decoded.InboundID)
		assert.Len(t, decoded.Items, 1)
		assert.Equal(t, "SKU001", decoded.Items[0].SKU)
	})

	t.Run("flagged event serialization", func(t *testing.T) {
		evt := FlaggedEvent{
			InboundID:  "inb_3",
			CustomerID: "cust_2",
			Reason:     "damaged goods",
			Timestamp:  time.Now(),
		}

		data, err := json.Marshal(evt)
		require.NoError(t, err)

		var decoded FlaggedEvent
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "inb_3", decoded.InboundID)
		assert.Equal(t, "damaged goods", decoded.Reason)
	})
}

func TestEventPublisher_Topics(t *testing.T) {
	assert.Equal(t, "wms.inbound.submitted", TopicInboundSubmitted)
	assert.Equal(t, "wms.inbound.approved", TopicInboundApproved)
	assert.Equal(t, "wms.inbound.flagged", TopicInboundFlagged)
}
