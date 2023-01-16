default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 TF_ACC_LOG=TRACE go test ./... -v $(TESTARGS) -timeout 120m
