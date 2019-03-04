#FROM gitlab/gitlab-runner:latest
#USER root
#RUN apt-get update -y
#RUN apt-get install -y bash curl jq


#COPY ./docker-entrypoint.sh /usr/bin/docker-entrypoint.sh
#RUN chmod +x /usr/bin/docker-entrypoint.sh /setup-gitlab-runner.sh
#ADD ./docker-entrypoint.sh /usr/bin/docker-entrypoint.sh
#ADD ./setup-gitlab-runner.sh /usr/bin/setup-gitlab-runner.sh
#RUN chmod a+x /usr/bin/docker-entrypoint.sh /usr/bin/setup-gitlab-runner.sh
#ADD docker-entrypoint.sh /docker-entrypoint.sh
#ADD setup-gitlab-runner.sh /setup-gitlab-runner.sh
#RUN chmod a+x /docker-entrypoint.sh /setup-gitlab-runner.sh
#CMD ["/docker-entrypoint.sh"]

#CMD ["sh", "-c", "/docker-entrypoint.sh"]
#CMD ["/docker-entrypoint.sh"]
#CMD ["docker-entrypoint.sh"]
#CMD ["/usr/bin/docker-entrypoint.sh"]
#CMD ["/bash", "/docker-entrypoint.sh"]
#CMD ["/bin/bash", "/docker-entrypoint.sh"]
#ENTRYPOINT [ "/docker-entrypoint.sh" ]
#CMD ["run", "--user=gitlab-runner", "--working-directory=/home/gitlab-runner", "&", "sh", "-c", "/docker-entrypoint.sh"]
#USER gitlab-runner

#CMD ["docker-entrypoint.sh"]

#-----------------------
#FROM ubuntu:14.04
#ADD https://github.com/Yelp/dumb-init/releases/download/v1.0.2/dumb-init_1.0.2_amd64 /usr/bin/dumb-init
#RUN chmod +x /usr/bin/dumb-init
#RUN apt-get update -y && \
#    apt-get upgrade -y && \
#    apt-get install -y ca-certificates wget apt-transport-https vim nano bash curl jq && \
#    apt-get clean && \
#    rm -rf /var/lib/apt/lists/*
#RUN echo "deb https://packages.gitlab.com/runner/gitlab-ci-multi-runner/ubuntu/ `lsb_release -cs` main" > /etc/apt/sources.list.d/runner_gitlab-ci-multi-runner.list && \
#    wget -q -O - https://packages.gitlab.com/gpg.key | apt-key add - && \
#    apt-get update -y && \
#    apt-get install -y gitlab-ci-multi-runner && \
#    wget -q https://github.com/docker/machine/releases/download/v0.7.0/docker-machine-Linux-x86_64 -O /usr/bin/docker-machine && \
#    chmod +x /usr/bin/docker-machine && \
#    apt-get clean && \
#    mkdir -p /etc/gitlab-runner/certs && \
#    chmod -R 700 /etc/gitlab-runner && \
#    rm -rf /var/lib/apt/lists/*
#ADD entrypoint /
#RUN chmod +x /entrypoint
#ADD docker-entrypoint.sh /docker-entrypoint.sh
#ADD setup-gitlab-runner.sh /setup-gitlab-runner.sh
#RUN chmod a+x /docker-entrypoint.sh /setup-gitlab-runner.sh
#VOLUME ["/etc/gitlab-runner", "/home/gitlab-runner"]
#ENTRYPOINT ["/usr/bin/dumb-init", "/docker-entrypoint.sh"]
#CMD ["run", "--user=gitlab-runner", "--working-directory=/home/gitlab-runner"]

#---------------------
FROM alpine

ADD https://github.com/Yelp/dumb-init/releases/download/v1.0.2/dumb-init_1.0.2_amd64 /usr/bin/dumb-init
RUN chmod +x /usr/bin/dumb-init

RUN apk add --update \
		bash \
		ca-certificates \
		git \
		openssl \
		wget \
        bash \
        curl \
        jq

RUN wget -O /usr/bin/gitlab-ci-multi-runner https://gitlab-ci-multi-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-ci-multi-runner-linux-amd64 && \
	chmod +x /usr/bin/gitlab-ci-multi-runner && \
	ln -s /usr/bin/gitlab-ci-multi-runner /usr/bin/gitlab-runner && \
	wget -q https://github.com/docker/machine/releases/download/v0.7.0/docker-machine-Linux-x86_64 -O /usr/bin/docker-machine && \
	chmod +x /usr/bin/docker-machine && \
	mkdir -p /etc/gitlab-runner/certs && \
	chmod -R 700 /etc/gitlab-runner

ADD entrypoint /
RUN chmod +x /entrypoint

ADD docker-entrypoint.sh /docker-entrypoint.sh
ADD setup-gitlab-runner.sh /setup-gitlab-runner.sh
RUN chmod a+x /docker-entrypoint.sh /setup-gitlab-runner.sh

VOLUME ["/etc/gitlab-runner", "/home/gitlab-runner"]
ENTRYPOINT ["/usr/bin/dumb-init", "/docker-entrypoint.sh"]
CMD ["run", "--user=gitlab-runner", "--working-directory=/home/gitlab-runner"]
