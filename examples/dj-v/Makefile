CS_BIN       ?= ../../cs
ANALYZER     ?= $(CURDIR)/scripts/analyze.py
ADAPTER_CMD  ?= deno run --allow-read --allow-write --allow-run adapter/src/index.ts

export CS_AUDIO_ANALYZER := $(ANALYZER)

.PHONY: all synth run build clean

all: synth run

build:
	cd ../.. && go build -o cs ./cmd/cs/

synth:
	$(CS_BIN) synth --verbose --output artifact.json

run:
	$(CS_BIN) run \
		--artifact artifact.json \
		--capabilities capabilities.json \
		--adapter "$(ADAPTER_CMD)" \
		--verbose

clean:
	rm -f artifact.json
	rm -rf output/
