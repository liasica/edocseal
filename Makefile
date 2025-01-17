.PHONY: dev, prod, generate-ent

dev:
	bash ./deploy.sh dev

prod:
	bash ./deploy.sh prod

generate-ent:
	bash ./generate-ent.sh
