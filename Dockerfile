FROM golang:1.19.0-alpine

RUN apk update &&  apk add git
WORKDIR /app
RUN go install github.com/cosmtrek/air@v1.29.0

# air -c [tomlファイル名] // 設定ファイルを指定してair実行(WORKDIRに.air.tomlを配置しておくこと)
CMD ["air", "-c", ".air.toml"]