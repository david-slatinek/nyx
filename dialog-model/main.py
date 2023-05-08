from concurrent import futures

import grpc
import torch
from transformers import AutoModelForCausalLM, AutoTokenizer

import schema_pb2 as pb2
import schema_pb2_grpc as pb2_grpc


class DialogServiceServicer(pb2_grpc.DialogServiceServicer):
    def __init__(self):
        self.chat_history_ids = torch.tensor([])
        self.bot_input_ids = torch.tensor([])
        self.step = 0
        self.dialog_id = 0

    def Dialog(self, request, context):
        if request.dialog_id != self.dialog_id:
            self.chat_history_ids = torch.tensor([])
            self.bot_input_ids = torch.tensor([])
            self.step = 0
            self.dialog_id = request.dialog_id

        question = request.text
        print(f"Question: {question}")

        input_ids = tokenizer.encode(question + tokenizer.eos_token, return_tensors="pt")
        self.chat_history_ids = torch.cat([self.chat_history_ids, input_ids], dim=-1)
        self.bot_input_ids = torch.cat([self.chat_history_ids, input_ids],
                                       dim=-1) if self.step > 0 else input_ids
        self.chat_history_ids = model.generate(self.bot_input_ids, max_length=1000, pad_token_id=tokenizer.eos_token_id)

        answer = tokenizer.decode(self.chat_history_ids[:, self.bot_input_ids.shape[-1]:][0], skip_special_tokens=True)

        print(f"Answer: {answer}")

        self.step += 1
        if self.step > 10000:
            self.chat_history_ids = torch.tensor([])
            self.bot_input_ids = torch.tensor([])
            self.step = 0
            self.dialog_id = 0

        return pb2.DialogResponse(answer=answer)


if __name__ == "__main__":
    tokenizer = AutoTokenizer.from_pretrained("microsoft/DialoGPT-large")
    model = AutoModelForCausalLM.from_pretrained("microsoft/DialoGPT-large")

    print("Model loaded")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    pb2_grpc.add_DialogServiceServicer_to_server(DialogServiceServicer(), server)
    server.add_insecure_port("[::]:9080")
    server.start()
    server.wait_for_termination()
