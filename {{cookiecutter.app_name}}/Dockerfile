# Build Stage
FROM {{cookiecutter.docker_build_image}}:{{cookiecutter.docker_build_image_version}} AS build-stage

LABEL app="{{cookiecutter.app_name}}"
LABEL REPO="https://{{cookiecutter.source_path}}/{{cookiecutter.app_name}}"

ARG GOPRIVATE={{cookiecutter.goprivate}}
ENV GOPRIVATE=${GOPRIVATE}

# --- Private module authentication (uncomment one) ---
# Prefer BuildKit secret mounts to avoid baking credentials into image layers.
# Option 1: SSH (use: DOCKER_BUILDKIT=1 docker build --ssh default .)
# RUN --mount=type=ssh git config --global url."git@github.com:".insteadOf "https://github.com/"
# Option 2: .netrc via BuildKit secret (use: docker build --secret id=netrc,src=$HOME/.netrc .)
# RUN --mount=type=secret,id=netrc,target=/root/.netrc true
# Option 3: Token via build arg (less secure — visible in image history)
# ARG GITHUB_TOKEN
# RUN echo "machine github.com login x-access-token password ${GITHUB_TOKEN}" > ~/.netrc && \
#     make build-alpine && rm -f ~/.netrc

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
