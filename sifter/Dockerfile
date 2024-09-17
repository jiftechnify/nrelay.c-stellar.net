FROM public.ecr.aws/docker/library/golang:1.21.3-alpine AS build
WORKDIR /go/src
COPY . /go/src/
RUN apk --update --no-cache add ca-certificates
RUN CGO_ENABLED=0 go build -o /go/bin/evsifter_whitelist -ldflags '-extldflags "-static"' cmd/whitelist/main.go
RUN CGO_ENABLED=0 go build -o /go/bin/evsifter_import_dm_wl -ldflags '-extldflags "-static"' cmd/import_dm_wl/main.go 

FROM dockurr/strfry:1.0.1
COPY --from=build /go/bin/* /app/plugin/
WORKDIR /app
EXPOSE 7777
ENTRYPOINT ["/app/strfry"]
CMD ["relay"]
