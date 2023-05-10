import os

import torch
from torch.utils.data import Dataset, DataLoader
from transformers import BartTokenizer, AdamW, BartForConditionalGeneration, TrainingArguments, Trainer


# Define the dataset
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
    # Initialize the tokenizer and model
    tokenizer = BartTokenizer.from_pretrained('facebook/bart-large')
    model = BartForConditionalGeneration.from_pretrained('facebook/bart-large')

    # Initialize the optimizer and set the learning rate
    optimizer = AdamW(model.parameters(), lr=1e-5)

    # Set the batch size and number of training epochs
    batch_size = 1
    num_epochs = 1

    # Set the file path and block size for the dataset
    block_size = 1

    # Initialize the dataset and data loader
    dataset = TextDataset(tokenizer, file_path="train.txt", block_size=block_size)
    train_loader = DataLoader(dataset, batch_size=batch_size, shuffle=True)

    # Train the model
    for epoch in range(num_epochs):
        for batch in train_loader:
            optimizer.zero_grad()
            loss = model(input_ids=batch, labels=batch)[0]
            loss.backward()
            optimizer.step()

    test_dataset = TextDataset(tokenizer, file_path="test.txt", block_size=block_size)
    test_loader = DataLoader(test_dataset, batch_size=batch_size, shuffle=True)

    eval_loss = evaluate(model, test_loader)
    print(f"Evaluation loss: {eval_loss:.4f}")
