FROM gitlab/gitlab-ce:latest

RUN apt-get update -y
RUN apt-get install -y jq
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
COPY ./setup-gitlab.sh /setup-gitlab.sh
COPY ./init-repo.sh /init-repo.sh
RUN chmod +x /docker-entrypoint.sh /setup-gitlab.sh /init-repo.sh


CMD ["/docker-entrypoint.sh"]
