#!/usr/bin/env python3
"""
MusicGen audio generation script for DJ-V.
Reads music parameters from stdin, generates audio using facebook/musicgen-small,
writes a WAV file, and outputs the file path as JSON.
"""

import json
import sys
import os
import time

def get_device():
    """Detect the best available device."""
    import torch
    if torch.backends.mps.is_available():
        return "mps"
    if torch.cuda.is_available():
        return "cuda"
    return "cpu"

def load_melody(melody_path: str, sample_rate: int):
    """Load a melody audio file for conditioning. Returns (waveform_tensor, sample_rate)."""
    import torch
    try:
        import torchaudio
        waveform, sr = torchaudio.load(melody_path)
        if sr != sample_rate:
            waveform = torchaudio.functional.resample(waveform, sr, sample_rate)
        if waveform.shape[0] > 1:
            waveform = waveform.mean(dim=0, keepdim=True)
        return waveform
    except Exception as e:
        print(f"Warning: failed to load melody from {melody_path}: {e}", file=sys.stderr)
        return None


def generate(params: dict, output_dir: str) -> str:
    """Generate audio from music parameters and return the output file path."""
    import torch
    from transformers import AutoProcessor, MusicgenForConditionalGeneration
    import scipy.io.wavfile

    device = get_device()
    print(f"Using device: {device}", file=sys.stderr)

    melody_path = params.get("melody_audio_path", "")
    use_melody = bool(melody_path and os.path.isfile(melody_path))

    model_name = "facebook/musicgen-melody" if use_melody else "facebook/musicgen-small"
    print(f"Loading MusicGen model ({model_name})...", file=sys.stderr)
    processor = AutoProcessor.from_pretrained(model_name)
    model = MusicgenForConditionalGeneration.from_pretrained(model_name)

    if device != "cpu":
        model = model.to(device)

    description = params.get("description", "A melodic electronic track")
    duration_seconds = 15

    max_tokens = int(duration_seconds * 50)

    print(f"Generating {duration_seconds}s of audio...", file=sys.stderr)
    print(f"Description: {description}", file=sys.stderr)

    start = time.time()

    if use_melody:
        print(f"Using melody conditioning from: {melody_path}", file=sys.stderr)
        melody_waveform = load_melody(melody_path, model.config.audio_encoder.sampling_rate)
        if melody_waveform is not None:
            inputs = processor(
                text=[description],
                audio=melody_waveform.squeeze().numpy(),
                sampling_rate=model.config.audio_encoder.sampling_rate,
                padding=True,
                return_tensors="pt",
            )
        else:
            use_melody = False
            inputs = processor(text=[description], padding=True, return_tensors="pt")
    else:
        inputs = processor(text=[description], padding=True, return_tensors="pt")

    if device != "cpu":
        inputs = {k: v.to(device) if hasattr(v, 'to') else v for k, v in inputs.items()}

    audio_values = model.generate(**inputs, max_new_tokens=max_tokens)
    elapsed = time.time() - start
    print(f"Generation took {elapsed:.1f}s", file=sys.stderr)

    audio_data = audio_values[0, 0].cpu().numpy()
    sample_rate = model.config.audio_encoder.sampling_rate

    os.makedirs(output_dir, exist_ok=True)
    output_path = os.path.join(output_dir, "dj-v-output.wav")
    scipy.io.wavfile.write(output_path, rate=sample_rate, data=audio_data)
    print(f"Wrote {output_path}", file=sys.stderr)

    return output_path

def main():
    input_data = json.loads(sys.stdin.read())

    # The input may be the music_parameters directly or nested under a key
    if isinstance(input_data, dict) and "music_parameters" in input_data:
        params = input_data["music_parameters"]
    elif isinstance(input_data, dict) and "description" in input_data:
        params = input_data
    else:
        params = input_data

    if isinstance(params, str):
        try:
            params = json.loads(params)
        except json.JSONDecodeError:
            params = {"description": params}

    script_dir = os.path.dirname(os.path.abspath(__file__))
    output_dir = os.path.join(script_dir, "..", "output")

    output_path = generate(params, output_dir)

    result = {
        "audio_file": output_path,
        "format": "wav",
        "duration_seconds": 15,
    }
    print(json.dumps(result))

if __name__ == "__main__":
    main()
