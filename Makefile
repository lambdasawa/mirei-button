.PHONY: \
	deploy \
	destroy \
	lint

deploy:
	cd infra &&\
		yarn cdk deploy &&\
		yarn post-deploy

destroy:
	cd infra && yarn cdk destroy

lint:
	lefthook install
	lefthook run pre-commit
