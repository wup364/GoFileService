FROM centos:centos8

WORKDIR /opts/fileservice

# RUN
RUN mkdir /opts/fileservice/conf \
 && mkdir /opts/fileservice/datas \
 && mkdir /opts/fileservice/.datas

# COPY
#COPY ./conf /opts/fileservice/conf
#COPY ./datas /opts/fileservice/datas
COPY ./webapps /opts/fileservice/webapps
# COPY ./app.ico /opts/fileservice/app.ico
COPY ./fileservice /opts/fileservice/fileservice

# CMD [  ]

# ENTRYPOINT []
ENTRYPOINT  ["./fileservice"]

# VOLUME[]
VOLUME ["/opts/fileservice/datas", "/opts/fileservice/.datas", "/opts/fileservice/conf"]
