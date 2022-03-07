# UI build
FROM node:14-buster AS ui
RUN apt-get update -y && apt-get install -y build-essential
RUN mkdir -p /home/app && chown -R node:node /home/app
WORKDIR /home/app
USER node
COPY --chown=node:node package*.json /home/app/
COPY --chown=node:node Makefile /home/app/
RUN make node_modules
COPY --chown=node:node ui /home/app/ui
RUN make ui

# Go build
FROM golang:1.17 AS go-build
# Expect to be passed these by running this via `make docker-gitops` so we don't copy all of .git/
ARG LDFLAGS="-X localbuild=true"
ARG GIT_COMMIT="_unset_"

# Add known_hosts entries for GitHub and GitLab
RUN mkdir ~/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN ssh-keyscan gitlab.com >> ~/.ssh/known_hosts

COPY Makefile /app/
WORKDIR /app
COPY go.* /app/
RUN go mod download
COPY --from=ui /home/app/cmd/gitops-server/cmd/dist/ /app/cmd/gitops-server/cmd/dist/
COPY . /app
# ignore the index.html dependency (which it otherwise would because node_modules is missing)
RUN LDFLAGS=$LDFLAGS GIT_COMMIT=$GIT_COMMIT make -o cmd/gitops-server/cmd/dist/index.html gitops-server

#  Distroless
FROM gcr.io/distroless/base as runtime
COPY --from=go-build /app/bin/gitops-server /gitops-server
COPY --from=go-build /root/.ssh/known_hosts /root/.ssh/known_hosts

ENTRYPOINT ["/gitops-server"]
