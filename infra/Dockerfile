FROM ubuntu:19.04

# Install tools
RUN apt update -yqq
RUN apt install -yqq supervisor git make

# Setup application directory
RUN mkdir /opt/ShellgeiBot
WORKDIR /opt/ShellgeiBot

# Add user
RUN useradd -m -d /home/ShellgeiBot ShellgeiBot

# Setup logs
RUN install -d -m 0755 -o ShellgeiBot -g ShellgeiBot /var/log/ShellgeiBot
RUN touch /var/log/ShellgeiBot/ShellgeiBot.log
RUN touch /var/log/ShellgeiBot/ShellgeiBot_error.log
RUN chown ShellgeiBot: /var/log/ShellgeiBot/ShellgeiBot*.log
