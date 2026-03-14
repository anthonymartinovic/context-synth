/**
 * DJ-V Adapter — receives ResolvedState from stdin, finds the generated
 * audio file, and plays it using the platform's built-in audio player.
 */

interface ResolvedState {
  artifact: unknown;
  outputs: Record<string, unknown>;
}

interface AudioOutput {
  audio_file: string;
  format: string;
  duration_seconds: number;
}

async function readStdin(): Promise<string> {
  const decoder = new TextDecoder();
  const chunks: string[] = [];

  for await (const chunk of Deno.stdin.readable) {
    chunks.push(decoder.decode(chunk));
  }

  return chunks.join("");
}

function getPlayCommand(
  filePath: string
): { cmd: string; args: string[] } {
  const os = Deno.build.os;
  switch (os) {
    case "darwin":
      return { cmd: "afplay", args: [filePath] };
    case "windows":
      return {
        cmd: "powershell",
        args: [
          "-c",
          `(New-Object Media.SoundPlayer '${filePath}').PlaySync()`,
        ],
      };
    case "linux":
      return { cmd: "aplay", args: [filePath] };
    default:
      throw new Error(`Unsupported OS: ${os}`);
  }
}

async function main() {
  const input = await readStdin();
  const state: ResolvedState = JSON.parse(input);

  // Find the audio file from generate_audio output
  let audioOutput: AudioOutput | null = null;

  const generateAudio = state.outputs["generate_audio"];
  if (generateAudio && typeof generateAudio === "object") {
    audioOutput = generateAudio as AudioOutput;
  }

  const audioFile = state.outputs["audio_file"];
  if (!audioOutput && audioFile && typeof audioFile === "object") {
    audioOutput = audioFile as AudioOutput;
  }

  if (!audioOutput?.audio_file) {
    console.error("No audio file found in resolved state");
    Deno.exit(1);
  }

  console.error(`Playing: ${audioOutput.audio_file}`);
  console.error(
    `Format: ${audioOutput.format}, Duration: ${audioOutput.duration_seconds}s`
  );

  const { cmd, args } = getPlayCommand(audioOutput.audio_file);

  const process = new Deno.Command(cmd, {
    args,
    stdout: "inherit",
    stderr: "inherit",
  });

  const result = await process.output();

  if (!result.success) {
    console.error(`Playback failed with code ${result.code}`);
    Deno.exit(1);
  }

  console.error("Playback complete");
}

main();
