.PHONY: \
	deploy \
	destroy

deploy:
	cd infra &&\
		yarn cdk deploy &&\
		yarn post-deploy

destroy:
	cd infra && yarn cdk destroy
