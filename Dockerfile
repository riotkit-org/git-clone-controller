FROM alpine:3.14 AS build

ADD .build/git-clone-operator /git-clone-operator
RUN chmod +x /git-clone-operator

# ---
FROM scratch AS run

COPY --from=build /git-clone-operator /usr/bin/git-clone-operator
RUN git-clone-operator check-binary

CMD ["git-clone-operator"]
