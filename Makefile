ENTRY_POINT = src/srch.go


run:
	go run $(ENTRY_POINT)
install:
	go install $(ENTRY_POINT)
