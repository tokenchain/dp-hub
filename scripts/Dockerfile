FROM golang:1.13

WORKDIR /go/src/github.com/cosmos/peggy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make install

RUN ixod init local --chain-id=darkpool-1x
RUN echo password | ixocli keys import validator ./deploy/test/test_validator_key
RUN ixod add-genesis-account $(ixocli keys show validator -a) 1000000000stake,1000000000atom
RUN echo password | ixod gentx --name validator
RUN ixod collect-gentxs