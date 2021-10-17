FROM openapitools/openapi-generator-cli:latest

RUN mkdir /api

ADD ../../api/api.yml /api/api.yml

RUN generate \
    -i /api/api.yml \
    -o /generated_ts_sdk \
    -g typescript-axios \
    --additional-properties=supportsES6=true

VOLUME [ "/generated_ts_sdk" ]