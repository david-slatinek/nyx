from concurrent import futures

import grpc
from transformers import BlenderbotTokenizer, BlenderbotForConditionalGeneration

import schema_pb2 as pb2
import schema_pb2_grpc as pb2_grpc


class DialogServiceServicer(pb2_grpc.DialogServiceServicer):
    def Dialog(self, request, context):
        question = request.text
        print(f"Question: {question}")

        inputs = tokenizer(question, return_tensors="pt")
        reply_ids = model.generate(**inputs, max_new_tokens=100)
        answer = tokenizer.decode(reply_ids[0], skip_special_tokens=True, clean_up_tokenization_spaces=False)

        print(f"Answer: {answer}")

        return pb2.DialogResponse(answer=answer)


if __name__ == "__main__":
    model_name = "facebook/blenderbot-400M-distill"
    tokenizer = BlenderbotTokenizer.from_pretrained(model_name)
    model = BlenderbotForConditionalGeneration.from_pretrained(model_name)

    print("Model loaded")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    pb2_grpc.add_DialogServiceServicer_to_server(DialogServiceServicer(), server)
    server.add_insecure_port("[::]:9080")
    server.start()
    server.wait_for_termination()
