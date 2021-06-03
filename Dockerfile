FROM golang:1.16.4-stretch

# Add new user with uid, pwd and guid
ENV USERNAME chainlink
ENV UID 1000
ENV GID 1000
ENV PW chainlink
ENV USER_HOME /home/${USERNAME}

RUN useradd --create-home --shell /bin/bash ${USERNAME} --uid=${UID} && \
    echo "${USERNAME}:${PW}" | chpasswd

# Copy files required for test execution
COPY go.mod ${USER_HOME}/go.mod
COPY go.sum ${USER_HOME}/go.sum
COPY main_test.go ${USER_HOME}/main_test.go
COPY Makefile ${USER_HOME}/Makefile
COPY contracts/ ${USER_HOME}/contracts/
COPY config/ ${USER_HOME}/config/

# Declare GOPATH, install deps and build test runner
RUN echo "export GOPATH=/go" >> ${USER_HOME}/.bashrc
RUN cd ${USER_HOME} && go mod download && go mod verify && go test -c -o test

RUN chown -R ${UID}:${GID} ${USER_HOME}

USER ${UID}:${GID}
WORKDIR /home/${USERNAME}

CMD ["bash", "-c", "${USER_HOME}/test -test.v"]
