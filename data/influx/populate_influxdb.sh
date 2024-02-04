#!/bin/bash

influx restore --token "${INFLUX_TOKEN}" --host "${INFLUX_HOST}" --full populate_data/