deployment-kubernetes/overlays/eu-test/config.${PROFILE}.yaml: config/config.${PROFILE}.yaml ## Copy config to overlays directory | PROFILE
	install -D config/config.${PROFILE}.yaml deployment-kubernetes/overlays/eu-test/config.${PROFILE}.yaml

kube.deploy:: deployment-kubernetes/overlays/eu-test/config.${PROFILE}.yaml ## Deploy with kustomize | PROFILE
	kubectl apply -k deployment-kubernetes/overlays/${PROFILE}