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
	    perl
	     
RUN gem install cureutils
RUN gem install matsuya

ENV PATH $PATH:/root/.egison/bin:/root/.cabal/bin
WORKDIR /

CMD ["bash"]
