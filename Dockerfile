FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
COPY . /go/src/github.com/opensourceways/xihe-aicc-finetune
WORKDIR /go/src/github.com/opensourceways/xihe-aicc-finetune
RUN GO111MODULE=on CGO_ENABLED=0 go build -o xihe-aicc-finetune -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"
RUN tar -xf ./infrastructure/aiccfinetuneimpl/tools/obsutil.tar.gz

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow tzdata git bash && \
    groupadd -g 5000 mindspore && \
    useradd -u 5000 -g mindspore -s /bin/bash -m mindspore

USER mindspore
WORKDIR /opt/app

COPY --chown=mindspore --from=BUILDER /go/src/github.com/opensourceways/xihe-aicc-finetune/xihe-aicc-finetune /opt/app
COPY --chown=mindspore --from=BUILDER /go/src/github.com/opensourceways/xihe-aicc-finetune/obsutil /opt/app
COPY --chown=mindspore --from=BUILDER /go/src/github.com/opensourceways/xihe-aicc-finetune/infrastructure/aiccfinetuneimpl/tools/upload_folder.sh /opt/app
RUN chmod 550 /opt/app/xihe-aicc-finetune
RUN chmod 550 /opt/app/obsutil
RUN chmod 550 /opt/app/upload_folder.sh

ENTRYPOINT ["/opt/app/xihe-aicc-finetune"]

