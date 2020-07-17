FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/xrootd-k8s-operator \
    USER_UID=1001 \
    USER_NAME=xrootd-operator

# install operator binary
COPY build/_output/bin/xrootd-k8s-operator ${OPERATOR}
COPY configmaps /configmaps

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
