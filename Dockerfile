FROM golang:1.16-alpine AS build
WORKDIR /go/src/project/
COPY . /go/src/project/
RUN go build -o /bin/deplist ./cmd/deplist/ 

FROM registry.fedoraproject.org/fedora:34
RUN dnf install -y \
    golang-bin-1.16 \
    yarnpkg \
    maven \
    rubygem-bundler \
    npm \
    && dnf clean all
COPY --from=build /bin/deplist /
ENTRYPOINT ["/deplist"]
