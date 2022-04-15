FROM alpine:3.14 AS build

ADD .build/git-clone-operator /git-clone-operator
RUN chmod +x /git-clone-operator

# ---
FROM gcr.io/distroless/static-debian11

COPY --from=build /git-clone-operator /usr/bin/git-clone-operator
RUN ["/usr/bin/git-clone-operator", "check-binary"]

ENTRYPOINT ["/usr/bin/git-clone-operator"]
