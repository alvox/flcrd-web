FROM postgres:11-alpine

#ENV POSTGRES_DB: flcrd
#ENV POSTGRES_USER: flcrd
#ENV POSTGRES_PASSWORD: flcrd

COPY schema.sql /docker-entrypoint-initdb.d
