# Neva

Flow Based Programming Environment

```shell
$ git clone git@github.com:emil14/neva.git

# backend
$ cd neva/cmd/server
$ go run .

# frontend
cd web
npm start
```

## Deps

- Docker
- Go
- NodeJS and NPM

## Generate SDK

### Go

```shell
docker run --rm
  -v ${PWD}:/app openapitools/openapi-generator-cli generate
  -i /app/api/api.yml \
  -g go-server  \
  -o /app/pkg/api
  --additional-properties=isGoSubmodule=false,packageName=sdk

mv pkg/api/go pkg/sdk
rm -r pkg/api
```

### TypeScript

```shell
docker run --rm
  -v ${PWD}:/app openapitools/openapi-generator-cli generate
  -i /app/api/api.yml \
  -g typescript \
  -o /app/web/sdk \
  --additional-properties=supportsES6=true
```


docker run --rm -v ${PWD}:/app openapitools/openapi-generator-cli generate -i /app/api/api.yml -g go-server -o /app/pkg/api --additional-properties=isGoSubmodule=false,packageName=sdk