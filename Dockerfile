# build stage
FROM golang:1.15.2-alpine AS build-env
ADD . /src
RUN cd /src && go build -o provy

# final stage
FROM alpine:3.12.1
WORKDIR /app
COPY --from=build-env /src/provy /usr/local/bin
ENTRYPOINT ["provy"]
