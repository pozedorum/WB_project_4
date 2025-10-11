all: mygrep

lint:
	golangci-lint run ./...
	go vet ./...

mygrep: clean
	go build -o mygrep ./cmd/main.go

clean:
	rm -rf mygrep tests/*.out *.out

test: mygrep
	@echo "=== Testing mygrep against standard grep ==="
	@echo ""
	
	@echo "1. Testing small file (test1.txt)..."
	@-./mygrep "test" tests/test1.txt > tests/mygrep_test1.out 2>&1 || true
	@-grep "test" tests/test1.txt > tests/grep_test1.out 2>&1 || true
	@if diff tests/mygrep_test1.out tests/grep_test1.out > /dev/null; then \
		echo "✓ test1.txt: PASS"; \
	else \
		echo "✗ test1.txt: FAIL"; \
		echo "Mygrep output:"; \
		cat tests/mygrep_test1.out; \
		echo "Grep output:"; \
		cat tests/grep_test1.out; \
	fi
	@echo ""
	
	@echo "2. Testing large file (test_large_file.txt) without concurrency..."
	@-./mygrep "test" tests/test_large_file.txt > mygrep_large.out 2>&1 || true
	@-grep "test" tests/test_large_file.txt > grep_large.out 2>&1 || true
	@if diff mygrep_large.out grep_large.out > /dev/null; then \
		echo "✓ test_large_file.txt (single): PASS"; \
	else \
		echo "✗ test_large_file.txt (single): FAIL"; \
		diff mygrep_large.out grep_large.out | head -20; \
	fi
	@echo ""
	
	@echo "3. Testing large file with concurrency (-Q 2)..."
	@-./mygrep -Q 2 "test" tests/test_large_file.txt > mygrep_concurrent.out 2>&1 || true
	@if diff mygrep_large.out mygrep_concurrent.out > /dev/null; then \
		echo "✓ test_large_file.txt (concurrent): PASS"; \
	else \
		echo "✗ test_large_file.txt (concurrent): FAIL"; \
		diff mygrep_large.out mygrep_concurrent.out | head -20; \
	fi
	@echo ""
	
	@echo "4. Testing case-insensitive search..."
	@-./mygrep -i "TEST" tests/test1.txt > mygrep_ci.out 2>&1 || true
	@-grep -i "TEST" tests/test1.txt > grep_ci.out 2>&1 || true
	@if diff mygrep_ci.out grep_ci.out > /dev/null; then \
		echo "✓ case-insensitive: PASS"; \
	else \
		echo "✗ case-insensitive: FAIL"; \
		diff mygrep_ci.out grep_ci.out; \
	fi
	@echo ""
	
	@echo "=== Test Summary ==="
	@echo "Check the *.out files for detailed results"

debug-test: mygrep
	@echo "=== Debug mode ==="
	@echo "Running mygrep directly:"
	@./mygrep "test" tests/test1.txt
	@echo "Exit code: $$?"