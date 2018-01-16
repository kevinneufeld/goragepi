FROM resin/raspberry-pi2-golang

ENV INITSYSTEM on

# RUN apt-get -q update && apt-get install -yq --no-install-recommends \
# 	build-essential \
# 	&& apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/kevinneufeld/GOragePi

COPY ./GOragePi ./app

# RUN go get -d -v ./...
# RUN go install -v ./...

# RUN go get GOragePi && go build

CMD ./app/GOragePi -pin $pin