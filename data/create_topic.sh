#!/bin/bash

kafka-topics --create --topic "$TOPIC_NAME" --replication-factor 1 --partitions 1 --bootstrap-server kafka:29093
echo "topic $TOPIC_NAME was created"