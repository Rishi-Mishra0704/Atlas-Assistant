import numpy as np
from piper import PiperVoice


def load_piper_voice(voice_path: str = "en_US-ryan-high.onnx") -> PiperVoice:
    return PiperVoice.load(voice_path)


def synthesize_speech(voice: PiperVoice, text: str) -> np.ndarray:
    audio_chunks = [
        np.frombuffer(chunk.audio_int16_bytes, dtype=np.int16)
        for chunk in voice.synthesize(text)
    ]

    if not audio_chunks:
        return np.array([], dtype=np.int16)

    return np.concatenate(audio_chunks)


def get_voice_sample_rate(voice: PiperVoice) -> int:
    return voice.config.sample_rate
