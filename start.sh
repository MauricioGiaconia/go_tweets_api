#!/bin/bash

# Iniciar Redis en segundo plano
redis-server &

# Iniciar API
/bin/api --db sqlite
