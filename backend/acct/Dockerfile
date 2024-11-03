FROM alpine:3.17

LABEL author="Yongjie Zhuang"
LABEL descrption="acct - simple personal accounting service"

RUN apk --no-cache add tzdata

WORKDIR /usr/src/

# binary is pre-compiled
COPY acct_build ./app_acct

ENV TZ=Asia/Shanghai

CMD ["./app_acct", "configFile=/usr/src/config/conf.yml"]
