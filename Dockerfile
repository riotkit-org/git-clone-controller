# syntax = docker/dockerfile:1.2
# ---
FROM golang:1.17 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
  go build -o bin/git-clone-operator .

# ---
FROM scratch AS run

COPY --from=build /work/bin/git-clone-operator /usr/local/bin/
RUN git-clone-operator check-binary

CMD ["git-clone-operator"]
