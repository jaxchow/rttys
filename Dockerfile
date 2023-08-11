FROM node:12 AS ui
WORKDIR /rttys-ui
COPY ui .
RUN npm install && npm run build

FROM golang:latest AS rttys
WORKDIR /rttys-build
COPY . .
COPY --from=ui /rttys-ui/dist ui/dist
RUN go env -w GOPROXY=https://goproxy.cn
RUN CGO_ENABLED=0 \
    VersionPath="rttys/version" \
    GitCommit=$(git log --pretty=format:"%h" -1) \
    BuildTime=$(date +%FT%T%z) \
    go build -ldflags="-s -w -X $VersionPath.gitCommit=$GitCommit -X $VersionPath.buildTime=$BuildTime"

FROM alpine:latest
WORKDIR /rttys
COPY --from=rttys /rttys-build/rttys /usr/bin/rttys
COPY --from=rttys /rttys-build/rttys.conf /rttys/rttys.conf
ENTRYPOINT ["/usr/bin/rttys","run","--conf","/rttys/rttys.conf"]
