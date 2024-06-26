FROM ubuntu
ENV DOCKER_DEFAULT_PLATFORM=linux/arm64

# Installing tools
RUN apt-get update && apt-get install -y git curl build-essential tree npm
RUN npm i -g n && n lts && npm i -g npm@latest
RUN npm install pm2 -g

# Installing Go
RUN curl -OL https://go.dev/dl/go1.22.4.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz && \
    rm go1.22.4.linux-amd64.tar.gz

# Adding Go to PATH
ENV GOPATH=/root/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

# Cloning repositories
WORKDIR /root
RUN git clone https://github.com/gnolang/gno.git && \ 
git clone https://github.com/gfant/shiken-api.git

# Installing GNO
RUN cd gno && make install

WORKDIR /
RUN mkdir apps

# Executing api
WORKDIR /root/shiken-api/executor
RUN go build main.go && chmod +x main && mv main /apps/
CMD ["pm2-runtime", "start", "/apps/main"]
EXPOSE 80
# docker build -t shiken .