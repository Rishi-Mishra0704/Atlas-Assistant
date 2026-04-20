import numpy as np
import sounddevice as sd


def record_audio(
    duration_seconds: float = 5,
    sample_rate: int = 16000,
    channels: int = 1,
) -> np.ndarray:
    audio = sd.rec(
        int(duration_seconds * sample_rate),
        samplerate=sample_rate,
        channels=channels,
        dtype="int16",
    )
    sd.wait()
    return audio


def audio_int16_to_float32(audio: np.ndarray) -> np.ndarray:
    return audio.flatten().astype(np.float32) / 32768.0


def play_audio(audio: np.ndarray, sample_rate: int) -> None:
    if audio.size == 0:
        return

    sd.play(audio, sample_rate)
    sd.wait()
