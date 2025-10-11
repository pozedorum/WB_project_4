all: clean mygrep

mygrep: clean
	go build -o mygrep ./cmd/main.go



clean:
	rm -rf mygrep

test: build
	@echo "=== Testing mygrep against standard grep ==="
	@echo ""
	
	@echo "1. Testing small file (test1.txt)..."
	@./mygrep "test" tests/test1.txt > mygrep_test1.out
	@grep "test" tests/test1.txt > grep_test1.out
	@if diff mygrep_test1.out grep_test1.out > /dev/null; then \
		echo "✓ test1.txt: PASS"; \
	else \
		echo "✗ test1.txt: FAIL"; \
		diff mygrep_test1.out grep_test1.out; \
	fi
	@echo ""
	
	@echo "2. Testing large file (test_large_file.txt) without concurrency..."
	@./mygrep "test" tests/test_large_file.txt > mygrep_large.out
	@grep "test" tests/test_large_file.txt > grep_large.out
	@if diff mygrep_large.out grep_large.out > /dev/null; then \
		echo "✓ test_large_file.txt (single): PASS"; \
	else \
		echo "✗ test_large_file.txt (single): FAIL"; \
		diff mygrep_large.out grep_large.out | head -20; \
	fi
	@echo ""
	
	@echo "3. Testing large file with concurrency (-Q 2)..."
	@./mygrep -Q 2 "test" tests/test_large_file.txt > mygrep_concurrent.out
	@if diff mygrep_large.out mygrep_concurrent.out > /dev/null; then \
		echo "✓ test_large_file.txt (concurrent): PASS"; \
	else \
		echo "✗ test_large_file.txt (concurrent): FAIL"; \
		diff mygrep_large.out mygrep_concurrent.out | head -20; \
	fi
	@echo ""
	
	@echo "4. Testing case-insensitive search..."
	@./mygrep -i "TEST" tests/test1.txt > mygrep_ci.out
	@grep -i "TEST" tests/test1.txt > grep_ci.out
	@if diff mygrep_ci.out grep_ci.out > /dev/null; then \
		echo "✓ case-insensitive: PASS"; \
	else \
		echo "✗ case-insensitive: FAIL"; \
		diff mygrep_ci.out grep_ci.out; \
	fi
	@echo ""
	
	@echo "=== Test Summary ==="
	@echo "Check the *.out files for detailed results"
	@echo "Clean up with: make clean"