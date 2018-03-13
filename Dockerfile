FROM golang:stretch

ENV DEBIAN_FRONTEND noninteractive

# Install aptly and required tools
RUN apt-get -q update                     \
    && apt-get -y install bash-completion \
                          bzip2           \
                          gnupg1          \
                          gpgv            \
                          graphviz        \
                          gpg             \
                          wget            \
                          xz-utils        \
                          gosu            \
                          ubuntu-archive-keyring \
    && echo "deb http://repo.aptly.info/ squeeze main" > /etc/apt/sources.list.d/aptly.list \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 9E3E53F19C7DE460 \
    && apt-get update \
    && apt-get -y install aptly \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY files/aptly.conf /etc/aptly.conf

VOLUME ["/var/lib/aptly"]
