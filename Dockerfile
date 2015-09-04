FROM scratch
ADD bin/scotty /scotty

EXPOSE 8080
ENTRYPOINT ["/scotty"]
