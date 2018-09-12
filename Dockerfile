FROM scratch
LABEL authors="Sangeet Kumar <sk@urantiatech.com>"
ADD cloudcms cloudcms
EXPOSE 8080
ENTRYPOINT ["/cloudcms", "--port=8080", "--dbFile=/master/cloudcms.db", "--indexDir=/master/cloudcms.index"]
