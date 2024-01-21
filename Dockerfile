FROM busybox:uclibc
WORKDIR /
COPY images-migrate-linux-amd64 /
CMD ["/images-migrate-linux-amd64"]