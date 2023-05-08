# Utilizzare un'immagine di base Golang
FROM golang:1.17

# Impostare la directory di lavoro all'interno del container
WORKDIR /app

# Copiare il file go.mod e go.sum nella directory di lavoro
COPY go.mod go.sum ./

# Scaricare tutte le dipendenze
RUN go mod download

# Copiare il codice sorgente nella directory di lavoro
COPY . .

# Compilare l'applicazione
RUN go build -o main .

# Utilizzare un'immagine di base leggera per l'esecuzione
FROM alpine:latest

# Installare le dipendenze necessarie
RUN apk --no-cache add ca-certificates libc6-compat

# Copiare il file binario dal primo stage
COPY --from=0 /app/main main

#ENV GIN_MODE=release

# Eseguire il binario
CMD ["./main"]
