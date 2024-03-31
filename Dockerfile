FROM alpine:latest

RUN ls -la
# Copiez votre application dans le conteneur
COPY README.md /app/README.md

# Définissez le répertoire de travail
WORKDIR /app

# Exécutez votre application lors du démarrage du conteneur
CMD ["./mon_api"]

# FROM alpine:latest

# # Copiez votre application dans le conteneur
# COPY mon_api_linux_amd64 /app/mon_api

# # Définissez le répertoire de travail
# WORKDIR /app

# # Exécutez votre application lors du démarrage du conteneur
# CMD ["./mon_api"]