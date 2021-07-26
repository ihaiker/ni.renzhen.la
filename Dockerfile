FROM python:3.9.6-alpine3.14
ADD entrypoint.sh /entrypoint.sh
RUN pip install mkdocs && \
    apk add go && \
    chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["build"]
