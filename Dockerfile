FROM golang:1.22
WORKDIR /src

COPY . .

RUN go build -o /mcfly main.go

ENV FLY_TOKEN ""
# ENV ADMIN_PASSWORD ""
# ENV DISCORD_TOKEN ""

ENTRYPOINT ["/mcfly"]
