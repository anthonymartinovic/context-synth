#!/usr/bin/env python3
"""
Audio feature extraction using librosa with HPSS.
Accepts a local audio file path as a CLI argument, extracts rich audio features,
and outputs a JSON object to stdout.

Usage: python3 analyze.py <path-to-audio-file>
"""

import json
import sys
import os

import librosa
import numpy as np


def estimate_key_mode(y_harmonic: np.ndarray, sr: int) -> tuple[int, int]:
    """Estimate musical key (0-11) and mode (0=minor, 1=major) from harmonic component."""
    chroma = librosa.feature.chroma_cqt(y=y_harmonic, sr=sr)
    chroma_mean = chroma.mean(axis=1)

    major_profile = np.array([6.35, 2.23, 3.48, 2.33, 4.38, 4.09, 2.52, 5.19, 2.39, 3.66, 2.29, 2.88])
    minor_profile = np.array([6.33, 2.68, 3.52, 5.38, 2.60, 3.53, 2.54, 4.75, 3.98, 2.69, 3.34, 3.17])

    best_key = 0
    best_mode = 0
    best_corr = -1.0

    for shift in range(12):
        shifted = np.roll(chroma_mean, shift)
        maj_corr = float(np.corrcoef(shifted, major_profile)[0, 1])
        min_corr = float(np.corrcoef(shifted, minor_profile)[0, 1])

        if maj_corr > best_corr:
            best_corr = maj_corr
            best_key = shift
            best_mode = 1
        if min_corr > best_corr:
            best_corr = min_corr
            best_key = shift
            best_mode = 0

    return best_key, best_mode


def analyze_audio(filepath: str) -> dict:
    y, sr = librosa.load(filepath, sr=None)

    y_harmonic, y_percussive = librosa.effects.hpss(y)

    # Tempo from percussive component for cleaner rhythm detection
    tempo, _ = librosa.beat.beat_track(y=y_percussive, sr=sr)
    if hasattr(tempo, '__len__'):
        tempo = float(tempo[0])
    else:
        tempo = float(tempo)

    # Key/mode from harmonic component for cleaner pitch detection
    key, mode = estimate_key_mode(y_harmonic, sr)

    # Energy (RMS of full signal)
    rms = librosa.feature.rms(y=y)
    energy = float(np.clip(np.mean(rms) * 5.0, 0.0, 1.0))

    # Danceability (onset strength of full signal)
    onset_env = librosa.onset.onset_strength(y=y, sr=sr)
    danceability = float(np.clip(np.std(onset_env) / 10.0, 0.0, 1.0))

    # Valence proxy (spectral centroid)
    spectral_centroid = librosa.feature.spectral_centroid(y=y, sr=sr)
    centroid_norm = np.mean(spectral_centroid) / (sr / 2)
    valence = float(np.clip(centroid_norm * 2.0, 0.0, 1.0))

    # Loudness
    loudness = float(20 * np.log10(np.mean(rms) + 1e-10))

    # Spectral contrast — peak-to-valley difference across 7 frequency bands
    contrast = librosa.feature.spectral_contrast(y=y, sr=sr)
    spectral_contrast = [round(float(v), 2) for v in contrast.mean(axis=1)]

    # Spectral rolloff — frequency below which most energy sits, normalized
    rolloff = librosa.feature.spectral_rolloff(y=y, sr=sr)
    spectral_rolloff = float(np.mean(rolloff) / (sr / 2))

    # Tonnetz — 6D tonal centroid from harmonic component
    tonnetz = librosa.feature.tonnetz(y=y_harmonic, sr=sr)
    tonnetz_features = [round(float(v), 4) for v in tonnetz.mean(axis=1)]

    # MFCCs — 13 mel-frequency cepstral coefficients as timbre fingerprint
    mfccs = librosa.feature.mfcc(y=y, sr=sr, n_mfcc=13)
    mfcc_features = [round(float(v), 2) for v in mfccs.mean(axis=1)]

    # Zero crossing rate — correlates with percussiveness/noisiness
    zcr = librosa.feature.zero_crossing_rate(y)
    zero_crossing_rate = round(float(np.mean(zcr)), 4)

    # Harmonic-to-percussive energy ratio
    rms_harmonic = float(np.mean(librosa.feature.rms(y=y_harmonic)))
    rms_percussive = float(np.mean(librosa.feature.rms(y=y_percussive)))
    total_rms = rms_harmonic + rms_percussive
    harmonic_ratio = round(rms_harmonic / total_rms, 2) if total_rms > 0 else 0.5

    return {
        "tempo": round(tempo, 1),
        "key": key,
        "mode": mode,
        "energy": round(energy, 2),
        "danceability": round(danceability, 2),
        "valence": round(valence, 2),
        "loudness": round(loudness, 1),
        "spectral_contrast": spectral_contrast,
        "spectral_rolloff": round(spectral_rolloff, 4),
        "tonnetz": tonnetz_features,
        "mfccs": mfcc_features,
        "zero_crossing_rate": zero_crossing_rate,
        "harmonic_ratio": harmonic_ratio,
    }


def main():
    if len(sys.argv) != 2:
        print("Usage: python3 analyze.py <audio-file-path>", file=sys.stderr)
        sys.exit(1)

    filepath = sys.argv[1]
    if not os.path.isfile(filepath):
        print(f"File not found: {filepath}", file=sys.stderr)
        sys.exit(1)

    print(f"Analyzing {os.path.basename(filepath)}...", file=sys.stderr)

    try:
        features = analyze_audio(filepath)
        print(json.dumps(features))
    except Exception as e:
        print(f"Analysis failed: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
