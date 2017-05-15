FROM scratch

ARG BINARY_NAME

COPY build/$BINARY_NAME /daemon

EXPOSE 10000 10001
CMD ["/daemon", "-daemon"]
HEALTHCHECK --interval=2m --timeout=3s CMD curl -f http://localhost:10000/healthz || exit 1
