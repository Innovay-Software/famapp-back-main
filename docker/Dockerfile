FROM golang:1.22.4

RUN apt-get update -y
RUN ln -fs /usr/share/zoneinfo/UTC /etc/localtime
RUN apt-get install -y vim
RUN apt-get install -y python3-launchpadlib
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:ubuntuhandbook1/ffmpeg6
RUN apt-get install -y ffmpeg
RUN apt-get install -y exif
RUN echo 'alias ll="ls -l"' >> ~/.bashrc
RUN go version
