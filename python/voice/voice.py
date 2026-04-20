import requests
from util.response import parse_api_response

from .piper_tts import get_voice_sample_rate, load_piper_voice, synthesize_speech
from .sound_device import audio_int16_to_float32, play_audio, record_audio
from .whisper_stt import load_whisper_model, transcribe_audio


def run_voice_loop():
    print("Loading Whisper...")
    stt = load_whisper_model(
        model_size="small",
        device="cpu",
        compute_type="int8",
    )

    print("Loading Piper voice...")
    tts = load_piper_voice("en_US-ryan-high.onnx")

    print("Atlas is listening. Press Ctrl+C to exit.\n")

    while True:
        try:
            print("Recording...")
            audio = record_audio(duration_seconds=5, sample_rate=16000)
            audio_float = audio_int16_to_float32(audio)

            print("Transcribing...")
            text = transcribe_audio(stt, audio_float, language="en").strip()

            if not text:
                print("Empty input, retrying.\n")
                continue

            print(f"You said: {text}")

            print("Calling Atlas API...")
            resp = requests.post(
                "http://localhost:8080/utterance",
                json={"text": text},
                timeout=30,
            )

            resp.raise_for_status()

            data = resp.json()
            payload = parse_api_response(data)

            reply = payload.get("reply", "")
            # closure = payload.get("closure", False)
            route = payload.get("route")
            print(f"Replying: {reply}")

            if reply:
                reply_audio = synthesize_speech(tts, reply)
                play_audio(reply_audio, get_voice_sample_rate(tts))

            if route == "close":
                print("Session closed by Atlas.")
                break

            print()

        except KeyboardInterrupt:
            print("\nExiting Atlas (manual interrupt).")
            break

        except Exception as e:
            print(f"Error: {e}\n")


def main() -> None:
    run_voice_loop()


if __name__ == "__main__":
    main()