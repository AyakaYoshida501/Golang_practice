FROM golang:1.18-alpine
# アップデートとgitとvimのインストール
RUN apk update && apk add git && apk add vim
# appディレクトリの作成
RUN mkdir /go/src/go_practice
# ワーキングディレクトリの設定
WORKDIR /go/src/go_practice
# ホストのファイルをコンテナの作業ディレクトリに移行
# イメージにファイルやフォルダを追加(COPYと役割は同じ？？どっちがいい？？)
COPY ./go_practice /go/src/go_practice
# ADD ./api-go /go/src/api-go
VOLUME /go/src/go_practice
EXPOSE 8080
# RUN go mod tidy
# GoApp起動
CMD ["go","run","main.go"]

