from datetime import datetime

import mlflow
import torch
from torch.utils.data import DataLoader, Dataset
from transformers import AdamW, BartForConditionalGeneration, BartTokenizer


class TextDataset(Dataset):
    def __init__(self, tokenizer, file_path, block_size):
        self.examples = []
        with open(file_path, encoding="utf-8") as f:
            text = f.read()
        tokenized_text = tokenizer.convert_tokens_to_ids(tokenizer.tokenize(text))
        for i in range(0, len(tokenized_text) - block_size + 1, block_size):
            self.examples.append(tokenizer.build_inputs_with_special_tokens(tokenized_text[i:i + block_size]))

    def __len__(self):
        return len(self.examples)

    def __getitem__(self, idx):
        return torch.tensor(self.examples[idx])


def evaluate(model, eval_loader):
    model.eval()
    total_loss = 0
    for batch in eval_loader:
        with torch.no_grad():
            loss = model(input_ids=batch, labels=batch)[0]
            total_loss += loss
    return total_loss / len(eval_loader)


if __name__ == "__main__":
    mlflow.sklearn.autolog()

    model_name = "facebook/bart-large"
    tokenizer = BartTokenizer.from_pretrained(model_name)
    model = BartForConditionalGeneration.from_pretrained(model_name)

    optimizer = AdamW(model.parameters())

    batch_size = 1
    num_epochs = 1
    block_size = 1

    dataset = TextDataset(tokenizer, file_path="../data/train.txt", block_size=block_size)
    train_loader = DataLoader(dataset, batch_size=batch_size, shuffle=True)

    start = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with mlflow.start_run(description=f"Training on {start}") as run:
        mlflow.log_params(model.config.__dict__)

        mlflow.log_param("batch_size", batch_size)
        mlflow.log_param("num_epochs", num_epochs)
        mlflow.log_param("block_size", block_size)

        mlflow.log_param("model_name", model_name)
        mlflow.log_param("optimizer", "AdamW")
        mlflow.set_tag("start", start)

        for epoch in range(num_epochs):
            for batch in train_loader:
                optimizer.zero_grad()
                loss = model(input_ids=batch, labels=batch)[0]
                loss.backward()
                optimizer.step()

        test_dataset = TextDataset(tokenizer, file_path="../data/test.txt", block_size=block_size)
        test_loader = DataLoader(test_dataset, batch_size=batch_size, shuffle=True)

        eval_loss = evaluate(model, test_loader)
        print(f"Evaluation loss: {eval_loss:.4f}")

        mlflow.log_metric("eval_loss", eval_loss)

        mlflow.set_tag("end", datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
