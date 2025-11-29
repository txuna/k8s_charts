.PHONY: deploy-prometheus
deploy-prometheus:
	@helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
		-n monitoring --create-namespace \
		-f ./values/prometheus/values.$(STAGE).yaml

.PHONY: remove-prometheus
remove-prometheus:
	@helm uninstall prometheus -n monitoring

.PHONY: deploy-nats
deploy-nats:
	@helm upgrade --install nats nats/nats \
		-n nats --create-namespace \
		-f ./values/nats/values.$(STAGE).yaml

.PHONY: remove-nats
remove-nats:
	@helm uninstall nats -n nats