FROM ubuntu:latest

RUN apt-get update -y && apt-get install -y \
	    ruby \
	    jq \
	    vim \
	    python3-dev \
	    python-dev \
	    nkf \
	    rs \
	    language-pack-ja \
	    pwgen \
	    bc \
	    perl \
	    toilet \
	    figlet
	     
RUN gem install cureutils
RUN gem install matsuya


ENV LANG ja_JP.UTF-8
WORKDIR /

CMD ["bash"]
