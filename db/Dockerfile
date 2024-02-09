FROM postgres:alpine3.18

ENV POSTGRES_PASSWORD=tameimpala \
    POSTGRES_USER=dbuser \
    POSTGRES_DB=wbl0

COPY ./db.sql /docker-entrypoint-initdb.d/db.sql

EXPOSE 5432