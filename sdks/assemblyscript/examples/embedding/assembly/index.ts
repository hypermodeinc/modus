import { models } from "@hypermode/functions-as";
import { EmbeddingsModel } from "@hypermode/models-as/models/experimental/embeddings";
import { OpenAIEmbeddingsModel } from "@hypermode/models-as/models/openai/embeddings";

// In this example, we will create embedding vectors from input text strings.
// For comparison, we'll do this with two different models.

export function testEmbeddingsWithMiniLM(texts: string[]): f32[][] {
  // In this example, we will use the MiniLM model to create embeddings.
  // See https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2 for more details about this model.
  // The EmbeddingsModel interface we use here is experimental and may change in the future.

  const model = models.getModel<EmbeddingsModel>("minilm");
  const input = model.createInput(texts);
  const output = model.invoke(input);
  return output.predictions;
}

export function testEmbeddingsWithOpenAI(texts: string[]): f32[][] {
  // See https://platform.openai.com/docs/api-reference/embeddings for more details
  // about the options available on the model, which you can set on the input object.

  const model = models.getModel<OpenAIEmbeddingsModel>("openai-embeddings");
  const input = model.createInput(texts);

  // For example, we can reduce the dimensions of the output vector if we want.
  // input.dimensions = 128;

  const output = model.invoke(input);

  // we can map the output to a 2D array of all the embeddings
  return output.data.map<f32[]>((d) => d.embedding);
}
