# Étape de lancement
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main
FROM alpine:latest

# COPY main /app/main
# Ajouter les fichiers binaires
# Copier tout le répertoire courant dans /app

COPY . /app
RUN chmod +x /app/amd64

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Exécuter l'application
CMD ["./amd64"]