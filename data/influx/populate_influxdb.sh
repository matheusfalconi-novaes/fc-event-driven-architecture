#!/bin/bash

influx restore --token "${INFLUX_ADMIN_TOKEN}" --host "${INFLUX_HOST}" --full populate_data/