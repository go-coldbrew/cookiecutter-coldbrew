# Build Stage
FROM {{cookiecutter.docker_build_image}}:{{cookiecutter.docker_build_image_version}} AS build-stage

LABEL app="build-{{cookiecutter.app_name}}"
LABEL REPO="https://{{cookiecutter.source_path}}/{{cookiecutter.app_name}}"

WORKDIR /app

COPY . .

RUN make build-alpine

# Final Stage
FROM {{cookiecutter.docker_image}}

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://{{cookiecutter.source_path}}/{{cookiecutter.app_name}}"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# add tz data and ca certs for HTTPS
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /opt/{{cookiecutter.app_name}}/bin

COPY --from=build-stage /app/bin/{{cookiecutter.app_name}} /opt/{{cookiecutter.app_name}}/bin/
RUN chmod +x /opt/{{cookiecutter.app_name}}/bin/{{cookiecutter.app_name}}

# Create appuser
RUN adduser -D -g '' {{cookiecutter.app_name}}
USER {{cookiecutter.app_name}}

EXPOSE 9090 9091

ENTRYPOINT ["/opt/{{cookiecutter.app_name}}/bin/{{cookiecutter.app_name}}"]
