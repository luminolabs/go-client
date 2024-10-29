# Using older python for torch compatibility,
# especially for multi-GPU training
FROM python:3.10-bullseye

# Install essentials
RUN apt update \
	&& apt install -y \
		build-essential \
		ca-certificates \
		curl \
		jq \
		git \
		libssl-dev \
		software-properties-common

# Install go 1.22.4
RUN curl -LO https://golang.org/dl/go1.22.4.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz \
    && rm go1.22.4.linux-amd64.tar.gz \
    && cp /usr/local/go/bin/go /usr/local/bin

# Set go environment variables
ENV GOROOT=/usr/local/go

# Install ethereum client 1.14.11
RUN curl -LO https://gethstore.blob.core.windows.net/builds/geth-alltools-linux-amd64-1.14.11-f3c696fa.tar.gz \
    && tar -C /usr/local -xzf geth-alltools-linux-amd64-1.14.11-f3c696fa.tar.gz \
    && rm geth-alltools-linux-amd64-1.14.11-f3c696fa.tar.gz \
    && cp /usr/local/geth-alltools-linux-amd64-1.14.11-f3c696fa/* /usr/local/bin

############################
### Install pipeline-zen ###
############################

# Upgrade pip
RUN python -m pip install --upgrade pip

# Work in this folder
WORKDIR /pipeline-zen-jobs

# Install these python libs outside of requirements.txt since they are large libraries
# and we don't want them to be build every time we add a new entry in requirements.txt
RUN pip install torch==2.4.1 transformers==4.44.2 datasets==3.0.0

# Install python libraries needed by the lib-celery
COPY pipeline-zen-src/lib-celery/requirements.txt ./requirements-lib-celery.txt
RUN pip install -r requirements-lib-celery.txt

# Install python libraries needed by the lib-common
COPY pipeline-zen-src/lib-common/requirements.txt ./requirements-lib-common.txt
RUN pip install -r requirements-lib-common.txt

# Install python libraries needed by the workflow
COPY pipeline-zen-src/lib-workflows/torchtunewrapper/requirements.txt ./requirements-workflow.txt
RUN pip install -r requirements-workflow.txt

# Copy scripts
COPY pipeline-zen-src/scripts scripts

# Copy application configuration folder
COPY pipeline-zen-src/app-configs app-configs

# Copy job configuration folder
COPY pipeline-zen-src/job-configs job-configs

# Copy celery workflow
COPY pipeline-zen-src/lib-celery lib-celery

# Copy lib-common source code
COPY pipeline-zen-src/lib-common lib-common

# Copy workflow source code
COPY pipeline-zen-src/lib-workflows lib-workflows

#########################
### Install go-client ###
#########################

WORKDIR /go-client

# Copy go-client source code
COPY . .

# Install go dependencies
RUN go mod tidy

## Generate bindings
RUN ./scripts/generate-bindings.sh

# Build go-client
RUN go build -o lumino