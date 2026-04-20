import numpy as np
from faster_whisper import WhisperModel


def load_whisper_model(
    model_size: str = "small",
    device: str = "cpu",
    compute_type: str = "int8",
) -> WhisperModel:
    return WhisperModel(model_size, device=device, compute_type=compute_type)


def transcribe_audio(
    model: WhisperModel,
    audio: np.ndarray,
    language: str = "en",
) -> str:
    segments, _ = model.transcribe(audio, language=language)
    return " ".join(segment.text.strip() for segment in segments).strip()
