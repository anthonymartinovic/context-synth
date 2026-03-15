# DJ-V — Context-Driven Music Generation

DJ-V is the first example application built on Context Synth. It turns governed context (documents + audio reference tracks) into generative music using Meta's MusicGen model.

## Prerequisites

### ollama + Llama 3.1

```bash
brew install ollama
ollama serve &
ollama pull llama3.1:8b
```

### Python + MusicGen dependencies

```bash
pip3 install -r scripts/requirements.txt
```

This installs MusicGen (transformers, torch, torchaudio, scipy) and librosa for audio analysis. The MusicGen model (`facebook/musicgen-small`, ~1.2GB) downloads automatically on first run.

### Audio reference tracks

Place MP3 files into `sources/audio/`. Each track is declared as an `audio:` source in `contextsynth.yml` with its own weight. During `cs synth`, librosa analyzes each track and extracts rich audio features (tempo, key, spectral characteristics, timbre fingerprint) that become part of the context artifact.

The example config ships with four reference tracks. To use your own, add MP3 files to `sources/audio/` and update `contextsynth.yml` to match:

```yaml
sources:
  - markdown: sources/md/vibe.md
    weight: 1.0
  - markdown: sources/md/constraints.md
    weight: 0.9
  - audio: sources/audio/your_track.mp3
    weight: 0.8
```

Copy the environment template and fill in the analyzer path:

```bash
cp ../../.env.example ../../.env
```

Edit `.env` so `CS_AUDIO_ANALYZER` points to the absolute path of the analysis script:

```
CS_AUDIO_ANALYZER=/full/path/to/examples/dj-v/scripts/analyze.py
```

If you use `make`, this is handled automatically — the Makefile sets `CS_AUDIO_ANALYZER` for you.

## Running DJ-V

The fastest way is `make`:

```bash
cd examples/dj-v
make          # synth + run in one shot
```

Or run each step individually:

```bash
make synth    # compile context artifact
make run      # resolve capabilities, generate + play audio
```

`make build` rebuilds the `cs` binary from source. `make clean` removes generated artifacts and output.

### Manual invocation

If you prefer not to use Make, the two steps are:

```bash
cd examples/dj-v
export CS_AUDIO_ANALYZER="$(pwd)/scripts/analyze.py"

cs synth --verbose --output artifact.json

cs run \
  --artifact artifact.json \
  --capabilities capabilities.json \
  --adapter "deno run --allow-read --allow-write --allow-run adapter/src/index.ts" \
  --verbose
```

This:
1. Loads the context artifact
2. Resolves `interpret_context` via Llama (derives music parameters)
3. Resolves `generate_audio` via MusicGen (generates 30s WAV file)
4. The adapter plays the audio

## How it works

```
contextsynth.yml
  │
  ├── sources/md/vibe.md (creative direction)
  ├── sources/md/constraints.md (musical constraints)
  └── sources/audio/*.mp3 (audio features via librosa + HPSS)
        │
        ▼
    cs synth → artifact.json
        │
        ▼
    cs run
        │
        ├── interpret_context (Llama)
        │   → music parameters + description
        │
        ├── generate_audio (MusicGen)
        │   → 30s WAV file
        │
        └── adapter (Deno)
            → audio playback
```

## Changing the output

- Edit `sources/md/vibe.md` to change the creative direction
- Edit `sources/md/constraints.md` to change musical constraints
- Edit `sources/md/composition.md` to describe the bar-by-bar structure
- Add or swap audio tracks in `sources/audio/` to change the musical reference
- Change source weights in `contextsynth.yml` to shift influence
- Re-run `cs synth` then `cs run` to hear the difference

## Contributing

DJ-V is an example application built on Context Synth — it uses the framework but is not part of it. Contributions that make DJ-V a better example are welcome. The framework itself lives in `internal/`, `integrations/`, and `cmd/` at the repo root.

### Extension points

**Richer context sources** — Add new markdown source documents (`sources/md/`) that describe genre, arrangement philosophy, or per-instrument direction. Each source gets its own weight in `contextsynth.yml`, so you control how much influence it has on the output.

**Melody conditioning** — `scripts/generate.py` already has a code path for melody-conditioned generation via `facebook/musicgen-melody`. Wire it by adding a `melody_audio_path` field to the music parameters produced by `interpret_context`, pointing at one of the reference tracks.

**Better prompt engineering** — The `interpret_context` prompt template in `capabilities.json` can be refined to produce more nuanced music parameters. The current template derives tempo, key, mode, energy, mood, density, and a description — there's room for instrumentation hints, dynamics, or arrangement instructions.

**Web-based adapter** — The current adapter plays audio via the system player (`afplay`/`aplay`). A browser-based adapter that renders a waveform visualization or streams playback would make the demo more shareable.

**Multi-track generation** — Generate separate stems (drums, bass, pads, leads) as individual capabilities in the graph, then mix them in the adapter. This would exercise deeper capability graphs with more dependency relationships.

**Output format options** — Support MP3/FLAC output, configurable duration, or sample rate via music parameters rather than hardcoded values in the generation script.

### What to keep in mind

- DJ-V does not import or modify anything in `internal/`. It communicates with the framework through the CLI (`cs synth`, `cs run`) and the protocol (artifact JSON, capability declarations, adapter interface).
- Changes to DJ-V should not require changes to the framework. If they do, that's a signal that something is missing from the framework's interfaces — raise it as an issue.
- The capability declarations in `capabilities.json` are DJ-V's, not the framework's. Add new capabilities, change prompts, swap executors — this is application-level configuration.
