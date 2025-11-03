include ../.envrc
# A simple Makefile to regenerate a specific gRPC proto file.

# Define the Python executable from the virtual environment, as requested.
VENV_PYTHON = venv/bin/python

# Phony target means this command does not produce a file named 'generate-proto'.
.PHONY: generate-proto

# This target executes the specific command you provided.
generate-proto:
	@echo "Generating Python gRPC code for job.proto..."
	$(VENV_PYTHON) -m grpc_tools.protoc -I protos --python_out=. --pyi_out=. --grpc_python_out=. protos/job.proto
	@echo "Done."

run:
	@echo "Running the server..."
	$(VENV_PYTHON) server.py
	@echo "Done."
