/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { models, http } from "@hypermode/modus-sdk-as";

import {
  OpenAIChatModel,
  DeveloperMessage,
  SystemMessage,
  UserMessage,
  TextContentPart,
  AudioContentPart,
  ImageContentPart,
  Image,
  Audio,
} from "@hypermode/modus-sdk-as/models/openai/chat";

// These examples demonstrate how to use audio or image data with OpenAI chat models.
// Currently, audio can be used for input or output, but images can be used only for input.

/**
 * This type is used in these examples to represent images or audio.
 */
class Media {
  // The content type of the media.
  contentType!: string;

  // The binary data of the media.
  // This value will be base64 encoded when used in an API response.
  data!: Uint8Array;

  // A text description or transcription of the media.
  text!: string;
}

/**
 * This function generates an audio response based on the instruction and prompt provided.
 */
export function generateAudio(instruction: string, prompt: string): Media {
  // Note, this is similar to the generateText example, but with audio output requested.

  // We'll generate the audio using an audio-enabled OpenAI chat model.
  const model = models.getModel<OpenAIChatModel>("audio-model");

  const input = model.createInput([
    new SystemMessage(instruction),
    new UserMessage(prompt),
  ]);

  input.temperature = 0.7;

  // Request audio output from the model.
  // Note, this is a convenience method that requests audio modality and sets the voice and format.
  // You can also set these values manually on the input object, if you prefer.
  input.requestAudioOutput("ash", "wav");

  const output = model.invoke(input);

  // Return the audio and its transcription.
  // Note that the message Content field will be empty for audio responses.
  // Instead, the text will be in the Message.Audio.Transcript field.
  const audio = output.choices[0].message.audio!;

  const media = <Media>{
    contentType: "audio/wav",
    data: audio.data,
    text: audio.transcript.trim(),
  };

  return media;
}

/**
 * This function generates text that describes the image at the provided url.
 * In this example the image url is passed to the model, and the model retrieves the image.
 */
export function describeImage(url: string): string {
  // Note that because the model retrieves the image, any URL can be used.
  // However, this means that there is a risk of sending data to an unauthorized host, if the URL is not hardcoded or sanitized.
  // See the describeRandomImage function below for a safer approach.

  const model = models.getModel<OpenAIChatModel>("text-generator");

  const input = model.createInput([
    UserMessage.fromParts([
      new TextContentPart("Describe this image."),
      new ImageContentPart(Image.fromURL(url)),
    ]),
  ]);

  const output = model.invoke(input);

  return output.choices[0].message.content.trim();
}

/**
 * This function fetches a random image, and then generates text that describes it.
 * In this example the image is retrieved by the function before passing it as data to the model.
 */
export function describeRandomImage(): Media {
  // Because this approach fetches the image directly, it is safer than the describeImage function above.
  // The host URL is allow-listed in the modus.json file, so we can trust the image source.

  // Fetch a random image from the Picsum API.  We'll just hardcode the size to make the demo simple to call.
  const response = http.fetch("https://picsum.photos/640/480");
  const data = Uint8Array.wrap(response.body);
  const contentType = response.headers.get("Content-Type")!;

  // Describe the image using the OpenAI chat model.
  const model = models.getModel<OpenAIChatModel>("text-generator");

  const input = model.createInput([
    UserMessage.fromParts([
      new TextContentPart("Describe this image."),
      new ImageContentPart(Image.fromData(data, contentType)),
    ]),
  ]);

  input.temperature = 0.7;

  const output = model.invoke(input);

  // Return the image and its generated description.
  const text = output.choices[0].message.content.trim();
  const media = <Media>{
    contentType,
    data,
    text,
  };

  return media;
}

/**
 * This function fetches a random "Harvard Sentences" speech file from OpenSpeech, and then generates a transcript from it.
 * The sentences are from https://www.cs.columbia.edu/~hgs/audio/harvard.html
 */
export function transcribeRandomSpeech(): Media {
  // Pick a random file number from the list of available here:
  // https://www.voiptroubleshooter.com/open_speech/american.html
  const numbers: i32[] = [
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 34, 35, 36, 37, 38, 39,
    40, 57, 58, 59, 60, 61,
  ];
  const num = numbers[<i32>(Math.random() * numbers.length)];

  // Fetch the speech file corresponding to the number.
  const url = `https://www.voiptroubleshooter.com/open_speech/american/OSR_us_000_00${num}_8k.wav`;
  const response = http.fetch(url);
  const data = Uint8Array.wrap(response.body);

  // Transcribe the audio using an audio-enabled OpenAI chat model.
  const model = models.getModel<OpenAIChatModel>("audio-model");

  const input = model.createInput([
    new DeveloperMessage(
      "Do not include any newlines or surrounding quotation marks in the response. Omit any explanation beyond the request.",
    ),
    UserMessage.fromParts([
      new TextContentPart(
        "Provide an exact transcription of the contents of this audio file.",
      ),
      new AudioContentPart(Audio.fromData(data, "wav")),
    ]),
  ]);

  const output = model.invoke(input);

  // Return the audio file and its transcript.
  const text = output.choices[0].message.content.trim();
  const media = <Media>{
    contentType: "audio/wav",
    data,
    text,
  };

  return media;
}
