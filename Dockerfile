FROM alpine:3.14 AS build

ADD .build/git-clone-controller /git-clone-controller
RUN chmod +x /git-clone-controller

# ---
FROM gcr.io/distroless/static-debian11

COPY --from=build /git-clone-controller /usr/bin/git-clone-controller
RUN ["/usr/bin/git-clone-controller", "check-binary"]

ENTRYPOINT ["/usr/bin/git-clone-controller"]
