start:
	go run cmd/kube-secret-sync/main.go start \
		--local --debug \
		--pod-namespace kube-secret-sync

generate:
	controller-gen object paths=./...