.PHONY: \
	deploy \
	destroy \
	run \
	lint

deploy:
	cd infra &&\
		yarn cdk deploy &&\
		yarn post-deploy

destroy:
	cd infra && yarn cdk destroy

run:
	docker-compose up

lint:
	lefthook install
	lefthook run pre-commit
