import torch
from transformers import AutoModelForCausalLM, AutoTokenizer

if __name__ == "__main__":
    tokenizer = AutoTokenizer.from_pretrained("microsoft/DialoGPT-large")
    model = AutoModelForCausalLM.from_pretrained("microsoft/DialoGPT-large")

    chat_history_ids = torch.tensor([])
    bot_input_ids = torch.tensor([])
    for step in range(5):
        question = input("Question: ")
        input_ids = tokenizer.encode(question + tokenizer.eos_token, return_tensors="pt")

        chat_history_ids = torch.cat([chat_history_ids, input_ids], dim=-1)

        bot_input_ids = torch.cat([chat_history_ids, input_ids], dim=-1) if step > 0 else input_ids

        chat_history_ids = model.generate(bot_input_ids, max_length=1000, pad_token_id=tokenizer.eos_token_id)

        answer = tokenizer.decode(chat_history_ids[:, bot_input_ids.shape[-1]:][0], skip_special_tokens=True)

        print(f"Answer: {answer}")
