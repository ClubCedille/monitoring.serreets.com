FROM node:16-alpine3.11 as build-node

LABEL maintainer="cedille"

# Installing make, g++, and python via apk will save headache on the M1 mac’s due to compatibility and architecture issues.
RUN apk --no-cache --virtual build-dependencies add \
        python \
        make \
        g++

WORKDIR /app

COPY assets/ .

RUN yarn install

# https://stackoverflow.com/questions/52103155/reading-an-environment-variable-in-react-which-was-set-by-docker
# This configuration makes available the variables passed via 
# ARG in docker-compose to Dockerfile during the image build 
# process and hence the variables can be in turn declared as env 
# variables which React can use during build process provided the webpack reads the env variables.
ARG REACT_APP_URL

ENV REACT_APP_URL=$REACT_APP_URL

RUN yarn build

FROM fitiavana07/nginx-react:latest

EXPOSE 80

COPY --from=build-node ./app/build /usr/share/nginx/html

