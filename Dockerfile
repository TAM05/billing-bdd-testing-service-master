FROM billing-bdd-testing-service-build as build
FROM alpine:3.6

ARG SERVICE

ENV APP=${SERVICE}

RUN apk add --no-cache ca-certificates && mkdir /app
COPY --from=build /${SERVICE} /app/${SERVICE}
COPY --from=build /go/src/github.com/utilitywarehouse/${SERVICE}/features /features

ENTRYPOINT ["/app/billing-bdd-testing-service"]
