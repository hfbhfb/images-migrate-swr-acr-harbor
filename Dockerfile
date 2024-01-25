FROM centos:centos7.9.2009
WORKDIR /
COPY images-migrate-linux-amd64 /
CMD ["/images-migrate-linux-amd64"]