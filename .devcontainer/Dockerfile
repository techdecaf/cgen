FROM golang:1.13.0-alpine

COPY .zshrc /tmp/.zshrc
RUN apk add git jq curl zsh \
    && sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" \
    && sh -c "$(curl -fsSL https://raw.github.com/techdecaf/tasks/master/install.sh)" \
    && tasks completion \
    && mv /tmp/.zshrc ~/.zshrc \
    && echo DONE