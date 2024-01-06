ENTRY_POINT = src/srch.go

run:
	go run $(ENTRY_POINT)

install:
	go install $(ENTRY_POINT)

quicktest:
	$(install)
	srch "ola" test/small 

testfile:
	$(instal)
	srch ipsum test/big

testdir:
	$(install)
	srch "ola" .

benchmarkfile:
	$(install)
	time srch ipsum test/big
	time grep ipsum test/big
