deploy:
	helm upgrade bootstrapper-testing-sample hack/testing/helm-sample \
		--install \
		--namespace bootstrapper-test \
		--create-namespace \
		--set image.repository=$(IMAGE) \
		--set image.tag=$(TAG) \
		--atomic

deploy/custom:
	helm upgrade bootstrapper-testing-sample hack/testing/helm-sample \
		--install \
		--namespace bootstrapper-test \
		--create-namespace \
		--set image.repository=$(IMAGE) \
		--set image.tag=$(TAG) \
		--atomic \
		-f hack/testing/helm-sample/_values.yaml

undeploy:
	helm uninstall bootstrapper-testing-sample \
		--namespace bootstrapper-test
