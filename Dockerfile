FROM ghcr.io/theshamuel/baseimg-go-build:1.20.1-1 as builder

ARG VER
ARG SKIP_TESTS
ENV GOFLAGS="-mod=vendor"

LABEL org.opencontainers.image.source https://github.com/theshamuel/medregistry20

RUN apk --no-cache add tzdata zip ca-certificates git

ADD . /build/medregestry20
ADD .golangci.yml /build/medregestry20/app/.golangci.yml
WORKDIR /build/medregestry20

#test
RUN \
    if [ -z "$SKIP_TESTS" ] ; then \
        go test -timeout=30s ./...; \
    else echo "[WARN] Skip tests" ; fi

#linter GolangCI
RUN \
    if [ -z "$SKIP_TESTS" ] ; then \
        golangci-lint run --skip-dirs vendor --config .golangci.yml ./...; \
    else echo "[WARN] Skip GolangCI linter" ; fi

RUN \
    ref=$(git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD); \
    version=${ref}_$(git log -1 --format=%h)_$(date +%Y%m%dT%H:%M:%S); \
    if [ -n "$VER" ] ; then \
    version=${VER}_${version}; fi; \
    echo "version=$version"; \
    go build -o medregestry20 -ldflags "-X main.version=${version} -s -w" ./app

FROM ghcr.io/theshamuel/baseimg-go-app:1.0-alpine3.17

WORKDIR /srv
COPY --from=builder /build/medregestry20/medregestry20 /srv/medregestry20

RUN chown -R appuser:appuser /srv
USER appuser

CMD [ "/srv/medregestry20", "server" ]