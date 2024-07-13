FROM golang:1.22-bookworm

WORKDIR /code

RUN git clone https://github.com/gjermundgaraba/mock-da.git
RUN cd mock-da && make build
RUN cp /code/mock-da/build/mock-da /usr/bin/mock-da

EXPOSE 7980