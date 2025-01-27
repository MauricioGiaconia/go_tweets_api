#!/bin/bash

# Iniciar Redis en segundo plano
redis-server &

# Esperar unos segundos para que Redis se inicie completamente
sleep 2

# Iniciar API
/bin/api --db sqlite
