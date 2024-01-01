ENTRY_POINT = src/srch.go


run:
	go run $(ENTRY_POINT)
install:
	go install $(ENTRY_POINT)
quicktest:
	go install $(ENTRY_POINT)
	srch "ola" test
