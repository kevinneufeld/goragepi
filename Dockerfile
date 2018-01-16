FROM resin/%%RESIN_MACHINE_NAME%%-golang:slim

ENV INITSYSTEM on

# RUN apt-get -q update && apt-get install -yq --no-install-recommends \
# 	build-essential \
# 	&& apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/GOragePi

COPY . ./

RUN go build

CMD ./GOragePi -pin $pin