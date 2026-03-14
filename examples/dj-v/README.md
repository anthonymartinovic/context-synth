# DJ-V — Context-Driven Music Generation

DJ-V is the first product built on Context Synth. It turns governed context (documents + audio reference tracks) into generative music using Meta's MusicGen model.

## Prerequisites

### ollama + Llama 3.1

```bash
brew install ollama
ollama serve &
ollama pull llama3.1:8b
```

### Python + MusicGen dependencies

```bash
pip3 install -r musicgen/requirements.txt
```

This installs MusicGen (transformers, torch, torchaudio, scipy) and librosa for audio analysis. The MusicGen model (`facebook/musicgen-small`, ~1.2GB) downloads automatically on first run.

### Audio reference tracks

Place MP3 files into `sources/tracks/`. Each track is declared as an `mp3:` source in `contextsynth.yml` with its own weight. During `cs synth`, librosa analyzes each track and extracts rich audio features (tempo, key, spectral characteristics, timbre fingerprint) that become part of the context artifact.

The example config ships with placeholder track names. Replace them with your own files and update `contextsynth.yml` to match:

```yaml
sources:
  - markdown: sources/vibe.md
    weight: 1.0
  - markdown: sources/constraints.md
    weight: 0.9
  - mp3: sources/tracks/your_track.mp3
    weight: 0.8
```

Set `CS_AUDIO_ANALYZER` to point to the analysis script (or add it to `.env`):

```bash
export CS_AUDIO_ANALYZER="/path/to/examples/dj-v/musicgen/analyze.py"
```

## Running DJ-V

### Step 1: Compile the context artifact

```bash
cd examples/dj-v
cs synth --no-llm
```

This compiles `sources/vibe.md`, `sources/constraints.md`, and audio track features into `artifact.json`.

With Llama for full LLM-backed extraction:

```bash
cs synth
```

### Step 2: Resolve capabilities and generate music

```bash
cs run \
  --artifact artifact.json \
  --capabilities capabilities.json \
  --adapter "deno run --allow-read --allow-write --allow-run adapter/src/index.ts" \
  --verbose
```

This:
1. Loads the context artifact
2. Resolves `interpret_context` via Llama (derives music parameters)
3. Resolves `generate_audio` via MusicGen (generates 15s WAV file)
4. The adapter plays the audio

## How it works

```
contextsynth.yml
  │
  ├── sources/vibe.md (creative direction)
  ├── sources/constraints.md (musical constraints)
  └── sources/tracks/*.mp3 (audio features via librosa + HPSS)
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
        │   → 15s WAV file
        │
        └── adapter (Deno)
            → audio playback
```

## Changing the output

- Edit `sources/vibe.md` to change the creative direction
- Edit `sources/constraints.md` to change musical constraints
- Add or swap audio tracks in `sources/tracks/` to change the musical reference
- Change source weights in `contextsynth.yml` to shift influence
- Re-run `cs synth` then `cs run` to hear the difference
