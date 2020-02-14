AWS_PROFILE := default
AWS_REGION := ap-northeast-1

# Project settings
PROJECT_NAME := hoge
ENV := dev
APP := ec2-chg-state

S3_BUCKET := ${PROJECT_NAME}-${ENV}-sam
STACK_NAME := ${PROJECT_NAME}-${ENV}-common-${APP}

NAME_TAG := "test-ec2"
WEBHOOK := "https://hooks.slack.com/services/hugahuga"
CHANNEL := "channel"

.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	-rm -rf ./build packaged.yaml
	-aws cloudformation delete-stack \
		--profile ${AWS_PROFILE} \
		--region ${AWS_REGION} \
		--stack-name ${STACK_NAME}

build:
	-mkdir ./build
	GOOS=linux GOARCH=amd64 go build -o build/handler ./handler

package: build
	sam package \
		--profile ${AWS_PROFILE} \
		--region ${AWS_REGION} \
		--s3-bucket ${S3_BUCKET} \
		--output-template-file packaged.yaml

deploy: package
	sam deploy \
		--profile ${AWS_PROFILE} \
		--region ${AWS_REGION} \
		--template-file packaged.yaml \
		--parameter-overrides \
			NameTag=${NAME_TAG} \
			Webhook=${WEBHOOK} \
			Channel=${CHANNEL} \
		--stack-name ${STACK_NAME} \
		--capabilities CAPABILITY_IAM
