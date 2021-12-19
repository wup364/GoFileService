FROM centos:centos8

WORKDIR /opts/fileservice

# RUN yum install golang

COPY ./conf /opts/fileservice/conf
COPY ./datas /opts/fileservice/datas
COPY ./webapps /opts/fileservice/webapps
COPY ./app.ico /opts/fileservice/app.ico
COPY ./fileservice /opts/fileservice/fileservice

# CMD [  ]

ENTRYPOINT  ["./fileservice"]

VOLUME ["/opts/fileservice/datas", "/opts/fileservice/conf"]