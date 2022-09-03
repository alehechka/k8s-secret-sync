start-regex:
	go run cmd/kube-secret-sync/main.go start \
		--exclude-regex-namespaces="kube[.]*" \
		--include-regex-namespaces="default[.]*" \
		--include-regex-secrets="test[.]*" \
		--exclude-regex-secrets "image-[.]*" \
		--local --debug \
		--secrets-namespace kube-secret-sync

generate:
	controller-gen object paths=./...