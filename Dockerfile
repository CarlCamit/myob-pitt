# Alpine 3.14 linux/amd64
FROM golang@sha256:ea5d6a7cf667df5041c09dc3741fc091fc07e6c8f996b580c1161d44313358b4 as builder

# Download OS dependencies
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create user
ENV USER=scoobydoo
ENV UID=29461

# https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /

COPY . .

# Build the binary with disabled symbol table and debugger
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/myob-pitt

###
FROM scratch

# Import dependencies and binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/myob-pitt /go/bin/myob-pitt

# Use unprivileged user
USER scoobydoo:scoobydoo

CMD ["/go/bin/myob-pitt"]