############################
# STEP 1 build executable binary
############################
# golang@1.14.0-alpine3.11 AMD64
FROM golang@sha256:e484434a085a28801e81089cc8bcec65bc990dd25a070e3dd6e04b19ceafaced as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
ENV USER=appuser UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
	--disabled-password     \
	--gecos ""              \
	--home "/nonexistent"   \
	--shell "/sbin/nologin" \
	--no-create-home        \
	--uid "${UID}"          \
	"${USER}"

WORKDIR $GOPATH/src/utekai/app/


# Fetch dependencies.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	# -ldflags='-w -s -extldflags "-static"' -a \
	-o /go/bin/utekai ./main

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/utekai /go/bin/utekai

# Use an unprivileged user
USER appuser:appuser

# Port on which the service will be exposed.
EXPOSE 9292

# Run the hello binary.
ENTRYPOINT ["/go/bin/utekai"]