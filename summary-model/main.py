from concurrent import futures

import grpc
from transformers import pipeline

import schema_pb2 as pb2
import schema_pb2_grpc as pb2_grpc


class SummaryServiceServicer(pb2_grpc.SummaryServiceServicer):
    def Summary(self, request, context):
        text = request.text
        print(f"Summary: {text}")

        summary = summarizer(text)

        print(f"Summary: {summary[0]['summary_text']}")

        return pb2.SummaryResponse(summary=summary[0]["summary_text"])


if __name__ == "__main__":
    summarizer = pipeline("summarization", model="philschmid/bart-large-cnn-samsum")

    print("Model loaded")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    pb2_grpc.add_SummaryServiceServicer_to_server(SummaryServiceServicer(), server)
    server.add_insecure_port("[::]:9050")
    server.start()
    server.wait_for_termination()
