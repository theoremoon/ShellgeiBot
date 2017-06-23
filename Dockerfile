FROM ubuntu

RUN apt-get update -y
RUN apt-get install -y ruby jq vim


RUN gem install cureutils

CMD ["bash"]
