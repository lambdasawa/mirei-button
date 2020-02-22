.PHONY: \
	deploy \
	destroy

deploy:
	cd infra && yarn cdk deploy

destroy:
	cd infra && yarn cdk destroy
