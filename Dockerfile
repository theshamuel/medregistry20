FROM theshamuel/baseimg-go-build:latest as builder

ARG VER
ARG SKIP_TESTS
ENV GOFLAGS="-mod=vendor"

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
    version="test"; \
    if [ -n "$VER" ] ; then \
    version=${VER}_$(date +%Y%m%d-%H:%M:%S); fi; \
    echo "version=$version"; \
    go build -mod=vendor  -o medregestry20 -ldflags "-X main.version=${version} -s -w" ./app

FROM theshamuel/baseimg-go-app:latest

WORKDIR /srv
COPY --from=builder /build/medregestry20/medregestry20 /srv/medregestry20

RUN chown -R appuser:appuser /srv
USER appuser

CMD [ "/srv/medregestry20", "server" ]