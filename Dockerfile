FROM ubuntu

RUN apt-get update -y
RUN apt-get install -y ruby jq vim python3-dev python-dev nkf rs language-pack-ja

RUN gem install cureutils

CMD ["bash"]
