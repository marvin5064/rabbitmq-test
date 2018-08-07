# Makefile
APPNAME  =`basename ${PWD}`

.PHONY: rabbitmq

run:
	go build
	./$(APPNAME)

rabbitmq:
	docker run -d \
		-p 127.0.0.1:8080:5672 \
		--name rabbitmq-management rabbitmq:management-alpine