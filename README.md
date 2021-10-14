- https://openapi-generator.tech/docs/generators/go-server

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
  -o /app/generated_go_sdk
  --additional-properties=isGoSubmodule=false,packageName=sdk

mv generated_go_sdk/go pkg/sdk
rm -r generated_go_sdk
```

docker run --rm -v ${PWD}:/app openapitools/openapi-generator-cli generate -i /app/api/api.yml -g go-server -o /app/generated --additional-properties=isGoSubmodule=false,packageName=sdk,featureCORS=true

### TypeScript

```shell
rm -r

docker run --rm
  -v ${PWD}:/app openapitools/openapi-generator-cli generate
  -i /app/api/api.yml \
  -g typescript \
  -o /app/generated_web \
  --additional-properties=supportsES6=true

cd generated_web
touch web/sdk
mv generated_web/dist web/sdk
rm -r generated_web
```

docker run --rm -v ${PWD}:/app openapitools/openapi-generator-cli generate -i /app/api/api.yml -g typescript -o /app/generated_web --additional-properties=supportsES6=true
