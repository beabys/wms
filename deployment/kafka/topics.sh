#!/bin/bash
# Kafka topic creation script for WMS
# Run after docker-compose up when Kafka is healthy.

set -e

BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-localhost:9092}"

echo "Creating topics on $BOOTSTRAP_SERVER..."

create_topic() {
  local topic=$1
  local partitions=${2:-1}
  local replication=${3:-1}

  docker exec wms-kafka \
    kafka-topics.sh \
    --create \
    --if-not-exists \
    --bootstrap-server "$BOOTSTRAP_SERVER" \
    --topic "$topic" \
    --partitions "$partitions" \
    --replication-factor "$replication"

  echo "  ✓ $topic"
}

# Customer events
create_topic "wms.customer.approved"
create_topic "wms.customer.suspended"

# Inbound events
create_topic "wms.inbound.submitted"
create_topic "wms.inbound.approved"
create_topic "wms.inbound.flagged"

echo ""
echo "All topics created successfully."
