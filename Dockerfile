FROM ubuntu:19.04

# set lang, locale
ENV LANG ja_JP.UTF-8
ENV TZ JST-9
ENV PATH /usr/games:$PATH

# apt install
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y ruby \
      ruby-dev\
      ccze\
      screen tmux\
      ttyrec\
      timidity abcmidi\
      r-base\
      boxes\
      ash yash\
      jq\
      vim emacs\
      python3-dev\
      python-dev\
      nkf\
      rs\
      language-pack-ja\
      pwgen\
      bc\
      perl\
      toilet\
      figlet\
      haskell-platform\
      git\
      build-essential\
      mecab libmecab-dev mecab-ipadic mecab-ipadic-utf8 python-mecab\
      wget curl npm\
      bsdgames fortunes cowsay fortunes-off cowsay-off\
      datamash\
      gawk\
      libxml2-utils\
      zsh\
      num-utils\
      apache2-utils\
      fish\
      cowsay\
      lolcat\
      nyancat\
      imagemagick\
      moreutils\
      strace\
      whiptail\
      pandoc\
      postgresql-common\
      postgresql-client-common\
      icu-devtools\
      tcsh\
      libskk-dev\
      libkkc-utils\
      morsegen\
      dc\
      telnet\
      python3-pip\
      python-pip\
      busybox\
      parallel\
      rename\
      mt-st\
      ffmpeg\
      kakasi\
      dateutils\
      fonts-ipafont fonts-vlgothic\
      inkscape gnuplot\
      qrencode\
      fonts-nanum fonts-symbola fonts-noto-color-emoji\
      sl\
      chromium-browser chromium-chromedriver nginx\
      screenfetch\
      mono-mcs mono-runtime\
      firefox\
      lua5.3 php7.2 php7.2-cli php7.2-common

# gem install
RUN gem install cureutils matsuya takarabako snacknomama rubipara marky_markov

# pip/pip3 install
RUN pip3 install yq faker sympy numpy scipy matplotlib xonsh pillow asciinema "https://github.com/megagonlabs/ginza/releases/download/v1.0.1/ja_ginza_nopn-1.0.1.tgz"
RUN pip install sympy numpy scipy matplotlib pillow

# install egzact && egison
RUN curl -OL --retry 3 https://git.io/egison-3.7.14.x86_64.deb && dpkg -i ./egison-3.7.14.x86_64.deb && rm ./egison-3.7.14.x86_64.deb
RUN curl -OL --retry 3 https://git.io/egzact-1.3.1.deb && dpkg -i ./egzact-1.3.1.deb && rm ./egzact-1.3.1.deb

# install node and faker-cli, chemi
RUN npm install -g faker-cli chemi

# home-commands (echo-sd)
RUN git clone https://github.com/fumiyas/home-commands.git
RUN rm /home-commands/README.md
ENV PATH /home-commands:$PATH

# j
RUN wget http://www.jsoftware.com/download/j807/install/j807_linux64_nonavx.tar.gz -O j.tar.gz && tar xvzf j.tar.gz && rm j.tar.gz
ENV PATH $PATH:/j64-807/bin

# trdsql (apply sql to csv)
RUN wget https://github.com/noborus/trdsql/releases/download/v0.5.0/trdsql_linux_amd64.zip && unzip trdsql_linux_amd64.zip
RUN rm trdsql_linux_amd64.zip
ENV PATH $PATH:/trdsql_linux_amd64

# openjdk11
RUN wget https://download.oracle.com/java/GA/jdk11/9/GPL/openjdk-11.0.2_linux-x64_bin.tar.gz -O openjdk11.tar.gz && tar xzf openjdk11.tar.gz
RUN rm openjdk11.tar.gz
ENV PATH $PATH:/jdk-11.0.2/bin

# super unko...
RUN curl -OL --retry 3 https://git.io/superunko.deb && dpkg -i superunko.deb && rm superunko.deb

# nameko.svg
RUN wget https://gist.githubusercontent.com/KeenS/6194e6ef1a151c9ea82536d5850b8bc7/raw/85af9ec757308b8ca4effdf24221f642cb34703b/nameko.svg

# go 1.12, sayhuuzoku, gron
RUN wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz -O go.tar.gz && tar xzf go.tar.gz -C /usr/local && rm go.tar.gz
ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /root/go
ENV PATH $PATH:/root/go/bin
RUN mkdir /root/go
RUN go get -u github.com/YuheiNakasaka/sayhuuzoku && ln -s /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db /
RUN go get -u github.com/tomnomnom/gron
RUN go get -u github.com/ericchiang/pup
RUN go get -u github.com/sugyan/ttyrec2gif
RUN go get -u github.com/xztaityozx/owari
RUN go get -u github.com/jiro4989/align
RUN go get -u github.com/jiro4989/taishoku
RUN go get -u github.com/jiro4989/textimg

# whitespace
RUN git clone https://github.com/hostilefork/whitespacers.git && cp /whitespacers/ruby/whitespace.rb /usr/local/bin/whitespace && chmod a+x /usr/local/bin/whitespace && rm -rf /whitespacers


# tukubai
RUN git clone https://github.com/usp-engineers-community/Open-usp-Tukubai
WORKDIR /Open-usp-Tukubai
RUN make install
WORKDIR /

# julia
RUN wget -O julia.tar.gz https://julialang-s3.julialang.org/bin/linux/x64/1.1/julia-1.1.0-linux-x86_64.tar.gz && tar xf julia.tar.gz && rm julia.tar.gz &&  ln -s $(realpath $(ls | grep -E "^julia") )/bin/julia /usr/local/bin/julia

# rust, rargs
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH /root/.cargo/bin:$PATH
RUN cargo install --git https://github.com/lotabout/rargs.git

# shellgei data
RUN git clone https://github.com/ryuichiueda/ShellGeiData

# imgout
RUN git clone https://github.com/ryuichiueda/ImageGeneratorForShBot
ENV PATH /ImageGeneratorForShBot:$PATH

# nignx
RUN /etc/init.d/nginx start

# zws, osquery, onefetch, sushiro, noc, bat
RUN wget https://raintrees.net/attachments/download/486/zws && chmod a+x ./zws
RUN wget https://pkg.osquery.io/deb/osquery_3.3.2_1.linux.amd64.deb -O osquery.deb && dpkg -i osquery.deb && rm osquery.deb
RUN wget https://github.com/o2sh/onefetch/releases/download/v1.5.2/onefetch_linux_x86-64.zip && unzip onefetch_linux_x86-64.zip && mv onefetch /usr/local/bin && rm onefetch_linux_x86-64.zip
RUN wget -nv https://raw.githubusercontent.com/redpeacock78/sushiro/master/sushiro && install -m 0755 sushiro /usr/local/bin/sushiro && rm sushiro && sushiro -f
RUN wget https://raw.githubusercontent.com/xztaityozx/noc/master/noc/noc/Program.cs && mcs Program.cs && rm Program.cs && mv Program.exe noc
RUN wget https://github.com/sharkdp/bat/releases/download/v0.10.0/bat_0.10.0_amd64.deb && sudo dpkg -i bat_0.10.0_amd64.deb && rm bat_0.10.0_amd64.deb

# echo-meme
RUN curl -OL --retry 3 https://git.io/echo-meme.deb && dpkg -i echo-meme.deb && rm echo-meme.deb

# bash5.0
RUN wget ftp://ftp.gnu.org/pub/gnu/bash/bash-5.0.tar.gz && tar xf bash-5.0.tar.gz && rm bash-5.0.tar.gz
WORKDIR bash-5.0
ENV CC cc
RUN ./configure && make && make install
WORKDIR /
RUN rm -rf bash-5.0

# awk5.0
RUN curl -OL --retry 3 https://ftp.gnu.org/gnu/gawk/gawk-5.0.0.tar.gz && tar -xpzf gawk-5.0.0.tar.gz && rm gawk-5.0.0.tar.gz
WORKDIR gawk-5.0.0
RUN ./configure --program-suffix="-5.0.0" && make && make install
WORKDIR /
RUN rm -rf gawk-5.0.0

# Support Japanese era name. Delete this line after new glibc (more than 2.29) is available
RUN curl -L --retry 3 "https://sourceware.org/git/?p=glibc.git;a=blob_plain;f=localedata/locales/ja_JP" > /usr/share/i18n/locales/ja_JP && localedef -i ja_JP -f UTF-8 ja_JP

RUN curl -O https://www.unicode.org/Public/UCD/latest/ucd/NormalizationTest.txt
RUN curl -O https://www.unicode.org/Public/UCD/latest/ucd/NamesList.txt

RUN CGO_LDFLAGS="`mecab-config --libs`" CGO_CFLAGS="-I`mecab-config --inc-dir`" go get -u github.com/ryuichiueda/ke2daira

CMD bash
