FROM golang:1.21 as builder 
ENV CGO_ENABLED 0
ARG VERSION
ARG ENVIRONMENT
ARG APIPORT

COPY . /coffee-api
WORKDIR /coffee-api/cmd/api

# Build the Coffee-Api, passing in VERSION from the Makefile 
RUN go build -ldflags="-X 'main.VERSION=${VERSION}' -X 'main.PORT=${APIPORT}' -X 'main.ENV=${ENVIRONMENT}'" -o coffee-api

WORKDIR /coffee-api/cmd/tooling

# Build Admin Migration Tool which will be used by the InitContainer
RUN go build -o migrations

FROM alpine:3.19
# Keep these ARGS in the final image
ARG BUILD_DATE
ARG BUILD_REF
ARG PORT
ARG DB_DSN
ENV DB_DSN=${DB_DSN}
# Set initial coffee-api admin password for db seeding.
# This is used by migrations binary.
ARG ADMIN_PASSWD
ENV ADMIN_PASSWD=${ADMIN_PASSWD}

# Ensure we have a valid user and group
RUN addgroup -g 1000 -S api-user && \
    adduser -u 1000 -G api-user -S api-user

# Copy application binary from builder image
COPY --from=builder --chown=api-user:api-user /coffee-api/cmd/api/coffee-api /cmsc/coffee-api
COPY --from=builder --chown=api-user:api-user /coffee-api/cmd/ui /cmsc/ui
COPY --from=builder --chown=api-user:api-user /coffee-api/cmd/tooling/migrations /cmsc/migrations

USER api-user
WORKDIR /cmsc
EXPOSE ${APIPORT}
CMD ["./coffee-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="coffee-api" \
      org.opencontainers.image.authors="Joshua Seals, Theodore Banta, Caleb Brennan" \
      org.opencontainers.image.source="https://github.com/tbanta5/CMSC" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="UMGC_CMSC_Group2_Productions"